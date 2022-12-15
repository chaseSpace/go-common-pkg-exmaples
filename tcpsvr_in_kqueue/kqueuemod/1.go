//go:build freebsd && darwin
// +build freebsd,darwin

package kqueuemod

import (
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_kqueue/socketmod"
	"log"
	"syscall"
)

type EventLoop struct {
	KqueueFileDescriptor int
	SocketFileDescriptor int
}

func NewEventLoop(s *socketmod.Socket) (*EventLoop, error) {
	// 创建了一个新的内核事件队列，待会儿用来订阅新socket连接的事件和轮训队列
	kQueue, err := syscall.Kqueue()
	if err != nil {
		return nil,
			fmt.Errorf("failed to create kqueue file descriptor (%v)", err)
	}
	// 构造一个event对象 传递给kqueue对象，表示我要订阅这个fd上的某些事件
	changeEvent := syscall.Kevent_t{
		Ident:  uint64(s.Fd),                       // 文件描述符标识
		Filter: syscall.EVFILT_READ,                // 处理事件的 Filter: 只关心传入连接的事件READ
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE, // 代表对这个事件要执行操作的 Flag：添加 & 启用
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	changeEventRegistered, err := syscall.Kevent(
		kQueue,
		[]syscall.Kevent_t{changeEvent},
		nil, nil,
	)
	if err != nil || changeEventRegistered == -1 {
		return nil,
			fmt.Errorf("failed to register change event (%v)", err)
	}

	return &EventLoop{
		KqueueFileDescriptor: kQueue,
		SocketFileDescriptor: s.Fd,
	}, nil
}

type Handler func(*socketmod.Socket)

func (eventLoop *EventLoop) Handle(handler Handler) {
	for {
		newEvents := make([]syscall.Kevent_t, 10)
		numNewEvents, err := syscall.Kevent( // BLOCK WAIT
			eventLoop.KqueueFileDescriptor,
			nil,
			newEvents,
			nil,
		)
		if err != nil {
			continue
		}
		log.Println("eventLoop new events...")

		for i := 0; i < numNewEvents; i++ {
			currentEvent := newEvents[i]
			eventFileDescriptor := int(currentEvent.Ident)
			// 处理 客户端关闭连接 事件
			if currentEvent.Flags&syscall.EV_EOF != 0 {
				// client closing connection
				syscall.Close(eventFileDescriptor)
				log.Println("event: close")
			} else if eventFileDescriptor == eventLoop.SocketFileDescriptor {
				// new incoming connection 新的客户端连接请求
				log.Println("event: new conn")
				socketConnection, _, err := syscall.Accept(eventFileDescriptor)
				if err != nil {
					log.Println("eventLoop Accept conn err:", err)
					continue
				}
				socketEvent := syscall.Kevent_t{
					Ident:  uint64(socketConnection),
					Filter: syscall.EVFILT_READ,
					Fflags: 0,
					Data:   0,
					Udata:  nil,
				}
				socketEventRegistered, err := syscall.Kevent(
					eventLoop.KqueueFileDescriptor,
					[]syscall.Kevent_t{socketEvent},
					nil,
					nil,
				)
				if err != nil || socketEventRegistered == -1 {
					log.Println("eventLoop register new conn err:", err)
					continue
				}
			} else if currentEvent.Filter&syscall.EVFILT_READ != 0 {
				// data available -> forward to handler
				// 某个客户端连接有数据进来了
				log.Println("event: new data")
				handler(&socketmod.Socket{
					Fd: eventFileDescriptor,
				})
			}
			// ignore all other events
		}
	}
}
