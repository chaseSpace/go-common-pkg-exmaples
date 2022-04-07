//go:build linux
// +build linux

package epoll

import (
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpcomm/conn"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpcomm/socket"
	"log"
	"sync"
	"syscall"
	"time"
)

// ET: EPOLL的边缘触发模式（edge-trigger），比LT模式复杂点，但性能高点
// 	-	读特性：除非有新数据进入socket read buffer，否则epoll_wait不会返回该socket的可读事件
// 	-	写特性：除非一次把socket write buffer写满，否则epoll_wait不会返回该socket的可写事件

// LT：水平触发模式（level-trigger）
//  -	读特性：只要socket read buffer有数据，epoll_wait就总会返回该socket的可读事件
// 	-	写特性：只要socket write buffer没写满，epoll_wait就总会返回该socket的可写事件

type EventLoop struct {
	epollFd int
	sock    *socket.Socket

	cmLock  sync.RWMutex
	connMap map[int]*conn.TcpConnET
}

func (e *EventLoop) Close() {
	_ = syscall.Close(e.epollFd)
	_ = e.sock.Close()
}

// CleanThread 一个单独线程来清理那些已关闭或内部发生ERR的conn
// 当然，这不是一个很好的实现，因为这里的锁会影响新连接的建立速度
func (e *EventLoop) CleanThread() {
	cleanup := func() {
		e.cmLock.Lock()
		defer e.cmLock.Unlock()
		for _, c := range e.connMap {
			if c.Err() != nil {
				log.Println("CleanThread: delete conn, fd", c.SockFd())
				delete(e.connMap, c.SockFd())
			}
		}
	}
	for {
		cleanup()
		time.Sleep(time.Second)
	}
}

func (e *EventLoop) safeReadTcpConn(fd int, op func(c *conn.TcpConnET) error) (err error) {
	e.cmLock.RLock()
	defer e.cmLock.RUnlock()
	tc, ok := e.connMap[fd]
	if ok {
		err = op(tc)
	}
	return
}

func (e *EventLoop) safeAddTcpConn(fd int, c *conn.TcpConnET) {
	e.cmLock.Lock()
	defer e.cmLock.Unlock()
	e.connMap[fd] = c
}

func (e *EventLoop) safeRemoveTcpConn(fd int) {
	e.cmLock.Lock()
	tcpConn, ok := e.connMap[fd]
	if ok {
		delete(e.connMap, fd)
		e.cmLock.Unlock()
		// 从epoll事件队列中注销对该socket事件的监听(必须先于 关闭socket的步骤)
		err := syscall.EpollCtl(
			e.epollFd,
			syscall.EPOLL_CTL_DEL,
			fd,
			nil,
		)
		if err != nil {
			log.Println("safeRemoveTcpConn: del fd err", err)
			return
		}
		tcpConn.Close()
		log.Printf("safeRemoveTcpConn: fd:%d OK!\n", fd)
	}
}

func NewEventLoop(ip string, port int) (et *EventLoop, err error) {
	sock, err := socket.Listen(ip, port)
	if err != nil {
		log.Println("Failed to create Socket:", err)
		return nil, err
	}
	defer func() {
		if et == nil {
			_ = sock.Close()
		}
	}()
	// 创建了一个新的内核事件队列，待会儿用来订阅新socket连接的事件
	// size用来告诉内核这个epoll实例同时监听的最大fd数目，但从linux内核2.6.8版本开始已弃用此参数，由内核自动分配，但必须大于0
	size := 1
	epollFd, err := syscall.EpollCreate(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create epoll file descriptor (%v)", err)
	}
	// 构造一个event对象 传递给epoll实例，表示我要订阅这个fd上的某些事件
	changeEvent := syscall.EpollEvent{
		Events: syscall.EPOLLIN | syscall.EPOLLET, // 订阅 IN（可读），并设置ET模式
		Fd:     int32(sock.Fd),
		Pad:    0,
	}

	err = syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, sock.Fd, &changeEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to register change event (%v)", err)
	}
	return &EventLoop{
		epollFd: epollFd,
		sock:    sock,
		connMap: make(map[int]*conn.TcpConnET),
	}, nil
}

func (e *EventLoop) Listen() {
	for {
		newEvents := make([]syscall.EpollEvent, 10) // 每次处理10个事件
		// 阻塞等待新的事件
		numNewEvents, err := syscall.EpollWait(
			e.epollFd, // epoll实例FD
			newEvents, // 待处理的事件数组结构，若有事件会填充到数组
			-1,        // 毫秒，表示在没有检测到事件发生时最多等待的时间, 负数则无限等待
		)
		if err == syscall.EINTR { // interrupted system call 系统中断，可忽略
			continue
		}
		if err != nil {
			log.Println("EpollWait err", err)
			continue
		}
		log.Printf("eventLoop new %d events ...\n", numNewEvents)

		for i := 0; i < numNewEvents; i++ {
			event := newEvents[i]
			eventFd := int(event.Fd)
			// 处理 ERR 事件（client关闭conn不会触发此事件）
			if event.Events&syscall.EPOLLERR != 0 && eventFd != e.sock.Fd {
				// client closing connection
				e.safeRemoveTcpConn(eventFd)
				log.Println("event: close")
			} else if eventFd == e.sock.Fd {
				// new incoming connection 新的客户端连接请求
				log.Println("event: new Conn")
				newSockFd, _, err := syscall.Accept(eventFd)
				if err != nil {
					log.Println("eventLoop Accept Conn err:", err)
					continue
				}
				// 设置socket非阻塞模式，以允许socket的read和write也是非阻塞的，这一步可选的，非阻塞模式可提高性能
				_ = syscall.SetNonblock(newSockFd, true)
				socketEvent := syscall.EpollEvent{
					/*
						因为现在是ET模式， 所以accept的时候需要一次性监听IN和OUT事件
						不得在read/write后修改事件，否则会重复收到readable/writeable事件 现采用如下逻辑
						-  accept后，监听 IN | ET | OUT
						-  read/write，不修改监听事件
					*/
					Events: syscall.EPOLLIN | syscall.EPOLLET | syscall.EPOLLOUT,
					Fd:     int32(newSockFd),
					Pad:    0,
				}
				// 监听新的conn socket
				syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_ADD, newSockFd, &socketEvent)
				c := conn.NewTcpConnET(socket.NewSocket(newSockFd))
				e.safeAddTcpConn(newSockFd, c)
				// 启动2个无限循环读写数据
				go c.ReadLoop()
				go c.WriteLoop()
			} else if event.Events&syscall.EPOLLIN != 0 {
				// ET模式下，这表示socket read buffer 刚从外部收到数据，收到一次触发一次
				log.Printf("event: Readable fd:%d\n", event.Fd)

			} else if event.Events&syscall.EPOLLOUT != 0 {
				// ET模式下，这表示socket write buffer刚从不可写变为可写（not writable->writeable）
				// 且刚启动时就会自动触发一次（client未发数据）
				log.Printf("event: Writeable fd:%d\n", event.Fd)
			}
		}
	}
}
