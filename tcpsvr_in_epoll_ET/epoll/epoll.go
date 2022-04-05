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
)

// ET: EPOLL的边缘触发模式（edge-trigger），比LT模式复杂点，但性能高点
// 	-	读特性：除非把socket read buffer读完，否则epoll_wait不会返回该socket的可读事件
// 	-	写特性：除非把socket write buffer写满，否则epoll_wait不会返回该socket的可写事件

// LT：水平触发模式（level-trigger）
//  -	读特性：只要socket read buffer有数据，epoll_wait就总会返回该socket的可读事件
// 	-	写特性：只要socket write buffer没写满，epoll_wait就总会返回该socket的可写事件

type EventLoop struct {
	epollFd int
	sock    *socket.Socket

	cmLock  sync.RWMutex
	connMap map[int]*conn.TcpConn
}

func (e *EventLoop) Close() {
	_ = syscall.Close(e.epollFd)
	_ = e.sock.Close()
}

func (e *EventLoop) safeReadTcpConn(fd int, op func(c *conn.TcpConn) error) (err error) {
	e.cmLock.RLock()
	defer e.cmLock.RUnlock()
	tc, ok := e.connMap[fd]
	if ok {
		err = op(tc)
	}
	return
}

func (e *EventLoop) safeAddTcpConn(fd int, c *conn.TcpConn) {
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
		Events: syscall.EPOLLIN | syscall.EPOLLET, // 订阅 IN（可读）和ERR事件，对于epoll实例的fd，只需要监听 IN（可读）事件，它不会有OUT（可写）事件
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
		connMap: make(map[int]*conn.TcpConn),
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
			// 处理 客户端关闭连接 事件
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
					Events: syscall.EPOLLIN | syscall.EPOLLET,
					Fd:     int32(newSockFd),
					Pad:    0,
				}
				// 监听新的conn socket
				syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_ADD, newSockFd, &socketEvent)
				c := conn.NewTcpConn(socket.NewSocket(newSockFd))
				e.safeAddTcpConn(newSockFd, c)
			} else if event.Events&syscall.EPOLLIN != 0 { // ET模式下，这表示 read buffer可读数据刚从0过渡到>0
				// 某个连接有数据进来了
				log.Printf("event: Readable fd:%d\n", event.Fd)
				err = e.safeReadTcpConn(eventFd, func(c *conn.TcpConn) error {
					c.Read()
					return c.Err()
				})
				if err != nil {
					e.safeRemoveTcpConn(eventFd)
				} else {
					// 修改监听的事件类型为：OUT（buffer可写） & ET模式
					event.Events = syscall.EPOLLOUT | syscall.EPOLLET
					err = syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_MOD, eventFd, &event)
				}
			} else if event.Events&syscall.EPOLLOUT != 0 { // ET模式下，表示write buffer可写空间刚从0过渡到>0
				log.Println("event: Writeable")
				err = e.safeReadTcpConn(eventFd, func(c *conn.TcpConn) error {
					c.WriteReply()
					return c.Err()
				})
				if err != nil {
					e.safeRemoveTcpConn(eventFd)
				} else {
					// 修改监听的事件类型为：IN（buffer可读） & ET模式
					event.Events = syscall.EPOLLIN | syscall.EPOLLET
					err = syscall.EpollCtl(e.epollFd, syscall.EPOLL_CTL_MOD, eventFd, &event)
				}
			}
		}
	}
}
