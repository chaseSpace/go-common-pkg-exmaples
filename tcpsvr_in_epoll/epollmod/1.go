//go:build linux
// +build linux

package epollmod

import (
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_epoll/socketmod"
	"log"
	"syscall"
)

const EPOLLET = 1 << 31

type EventLoop struct {
	epollFd int
	sock    *socketmod.Socket
}

func (e EventLoop) Close() {
	_ = syscall.Close(e.epollFd)
	_ = e.sock.Close()
}

func NewEventLoop(ip string, port int) (et *EventLoop, err error) {
	sock, err := socketmod.Listen(ip, port)
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
	// size用来告诉内核这个epoll实例同时监听的最大fd数目，但从linux内核2.6.8版本开始已弃用此参数，但必须大于0
	size := 1
	epollFd, err := syscall.EpollCreate(size)
	if err != nil {
		return nil, fmt.Errorf("failed to create epoll file descriptor (%v)", err)
	}
	// 构造一个event对象 传递给epollFd实例，表示我要订阅这个fd上的某些事件
	changeEvent := syscall.EpollEvent{
		Events: syscall.EPOLLIN | EPOLLET, // 订阅 IN（可读）和ERR事件
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
	}, nil
}

// Handler 其实现一般是开启新线/协程处理后续逻辑，才不会阻塞主线程epoll实例，最大化性能
// 另外，必须要控制最大并发线程数，避免耗尽内存，或在GC语言中造成过高延迟
type Handler func(*socketmod.Socket)

func (e *EventLoop) Handle(handler Handler) {
	for {
		newEvents := make([]syscall.EpollEvent, 10) // 每次处理10个事件
		// 阻塞等待新的事件
		numNewEvents, err := syscall.EpollWait(
			e.epollFd, // epoll实例FD
			newEvents, // 待处理的事件数组结构，若有事件会填充到数组
			-1,        // 10s 表示在没有检测到事件发生时最多等待的时间
		)
		if err != nil {
			continue
		}
		log.Printf("eventLoop new %d events ...\n", numNewEvents)

		for i := 0; i < numNewEvents; i++ {
			event := newEvents[i]
			eventFd := int(event.Fd)
			// 处理 客户端关闭连接 事件
			if event.Events&syscall.EPOLLERR != 0 {
				// client closing connection
				syscall.Close(eventFd)
				log.Println("event: close")
			} else if eventFd == e.sock.Fd {
				// new incoming connection 新的客户端连接请求
				log.Println("event: new conn")
				newSockFd, _, err := syscall.Accept(eventFd)
				if err != nil {
					log.Println("eventLoop Accept conn err:", err)
					continue
				}
				_ = syscall.SetNonblock(newSockFd, true)
				socketEvent := syscall.EpollEvent{
					Events: syscall.EPOLLIN | EPOLLET | syscall.EPOLLOUT, // 订阅 IN（可读）和ERR事件
					Fd:     int32(newSockFd),
					Pad:    0,
				}
				err = syscall.EpollCtl(
					e.epollFd,
					syscall.EPOLL_CTL_ADD,
					newSockFd,
					&socketEvent,
				)
				if err != nil {
					log.Println("eventLoop register new conn err:", err)
					continue
				}
			} else if event.Events&syscall.EPOLLIN != 0 {
				// data available -> forward to handler
				// 某个客户端连接有数据进来了
				log.Printf("event: new data fd:%d\n", event.Fd)
				go handler(&socketmod.Socket{ // 如果这里是异步，那必须设置 ET（边缘触发）模式：真的有数据进入socket时触发事件，而不是在缓冲区未读完时重复触发
					Fd: eventFd,
				})
			} else if event.Events&syscall.EPOLLOUT != 0 {
				// data available -> forward to handler
				// 某个客户端连接有数据进来了
				log.Println("event: write data ***")
				tmpSock := socketmod.Socket{Fd: eventFd}
				tmpSock.Write([]byte(`<html>
								  <head>
									<title>Epoll Response</title>
								  </head>
								  <body>%s</body>
								</html>`))
				event.Events = syscall.EPOLLIN | EPOLLET | syscall.EPOLLOUT
				syscall.EpollCtl(
					e.epollFd,
					syscall.EPOLL_CTL_MOD,
					int(event.Fd),
					&event,
				)
			}
		}
	}
}
