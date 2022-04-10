//go:build windows
// +build windows

package iocp

import (
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpcomm/socket"
	"golang.org/x/sys/windows"
	"log"
	"runtime"
	"unsafe"
)

type Iocp struct {
	sock     *socket.Socket
	hComPort windows.Handle
}

func (i *Iocp) Close() {
	err := windows.CloseHandle(i.hComPort)
	if err != nil {
		log.Println("CloseHandle err", err)
	}
	err = windows.WSACleanup() // 释放为我初始化的dll资源
	if err != nil {
		log.Println("WSACleanup err", err)
	}
	err = i.sock.Close()
	if err != nil {
		log.Println("sock.Close() err", err)
	}
}

func NewEventLoop(ip string, port int) (iocp *Iocp, err error) {
	sock, err := socket.Listen(ip, port)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			sock.Close()
		}
	}()
	// 指定我想要限制所使用的的win sockets最高版本，如果>=运行平台WinSock dll所支持的最低socket版本，WSAStartup()返回成功
	// 比如现在Ws2_32.dll的win sockets spec版本是2.2，支持2.2/2.1/2.0/1.1/1.0
	// 2.2在Windows Server 2008, Windows Vista, Windows Server 2003, Windows XP等较多系统上都支持
	// 对于Windows 95 and versions of Windows NT 3.51及更早的Windows系统，则1.1是最高支持的版本
	// （言下之意，若想编写的应用能够兼容更早期的win系统，应该使用较低的版本，比如1.1，但有得必有失，低版本可能失去某些性能特征）
	mainVer := int8(2)
	minorVer := int8(1)
	socketVer := MakeWord(mainVer, minorVer) // 第一参数表主版本，第二参数表副版本，这里表示我期望使用2.1版本

	// https://docs.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-wsastartup#remarks
	// 这个结构体用于接收实际使用的win socket版本信息
	wd := &windows.WSAData{
		Version:      0,           // 双字节数据，WSAStartup()返回成功后，此字段会被设置为平台dll<期望>我们使用的socket版本，即最终协商使用的版本（有可能并不等于传入的版本信息，所以下面有二次判断）
		HighVersion:  0,           // 双字节数据，WSAStartup()返回成功后，此字段会被设置为平台dll所支持的socket最高版本
		MaxSockets:   0,           // 能同时打开的最大socket数量，WinSock2及以后的协议栈已废弃此字段！
		MaxUdpDg:     0,           // 最大数据包size，WinSock2及以后的协议栈已废弃此字段！ ——可以在创建socket后通过getsockopt()查询SO_MAX_MSG_SIZE
		VendorInfo:   nil,         // 指明厂商特定的信息，WinSock2及以后的协议栈已废弃此字段！
		Description:  [257]byte{}, // WinSock 描述
		SystemStatus: [129]byte{}, // WinSock SystemStatus
	}

	// 首先必须调用的接口：使用wd初始化WSA
	// 内部逻辑：与运行的win平台dll层进行协商所要使用的socket版本信息
	err = windows.WSAStartup(uint32(socketVer), wd)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			windows.WSACleanup() // 释放初始化的dll资源
		}
	}()
	// NOTE:这里还可以在检查wd返回的版本是否符合需求，不符合可终止服务
	// >>这里的wd.Version也是同上面MakeWord()方式计算出来的，所以要反解出主副版本，再比对
	if LoByte(wd.Version) != uint16(mainVer) || HiByte(wd.Version) != uint16(minorVer) {
		err = fmt.Errorf("could not find a usable version of Winsock.dll")
		return
	}

	var maxSystemProcessIOThreadCnt = runtime.NumCPU()

	// 创建一个关联上述socket句柄的<完成端口>，类似epoll里面的事件队列概念
	// 这个端口将负责监听socket句柄上发生的所有I/O事件
	// 扩展：一个FileHandle只能绑定一个IOCP端口，但反之可以一绑多（通过多次调用此接口）
	hComPort, err := windows.CreateIoCompletionPort(
		sock.Fd,                             // FileHandle —— 一个支持overlapped I/O的端点，比如普通文件/socket句柄/命名管道
		windows.InvalidHandle,               // ExistingCompletionPort 已经存在的IOCP句柄
		0,                                   // CompletionKey 完成键，包含了指定I/O完成包的指定文件
		uint32(maxSystemProcessIOThreadCnt), // 内核层的真正并发同时执行最大线程数，若第二个参数为空则此参数被忽略，0也表示使用默认值（cpu核心数）
	)
	if err != nil {
		err = fmt.Errorf("could not find a usable version of Winsock.dll")
		return
	}
	iocp = &Iocp{sock: sock, hComPort: hComPort}

	// 这里的线程数其实就是限制并发调用GetQueuedCompletionStatus()的线程数
	// 如果超出传给CreateIoCompletionPort()的最后一个参数，对前者的调用就会阻塞
	for i := 0; i < maxSystemProcessIOThreadCnt; i++ {
		log.Println("NewEventLoop: start goroutine-", i)
		go iocp.workThread()
	}
	return iocp, nil
}

// 自定义结构体，用来保存client连接信息
type completionKey struct {
	sock *socket.Socket
	addr *windows.SockaddrInet4
}

type ioData struct {
	ov     *windows.Overlapped
	wsaBuf windows.WSABuf
	isRead bool
}

const BUF_SIZE = 10

// 工作线程不断的从windows内核接收I/O变化事件通知，并处理它们
func (i *Iocp) workThread() {
	var (
		bytesTrans uint32
		cpk        completionKey
		ovPlus     = &ioData{
			ov: &windows.Overlapped{},
		} // 之所以是双重指针，是因为内核可能设置它为nil
		timeout = uint32(windows.INFINITE) // 单次调用的等待时间：无限等待，也可设置一个数值，单位millisecond
	)

	// 启动一个无限循环
	for {
		// 从内核接收一个I/O完成事件，与epoll/kqueue的差别是，win内核已将IOCP端口绑定的端点上发生的IO数据copy到传入的结构体
		// 前者只是给个I/O ready事件，让我们自己去读写数据
		key := uintptr(unsafe.Pointer(&cpk))
		err := windows.GetQueuedCompletionStatus(i.hComPort, &bytesTrans, &key, &ovPlus.ov, timeout)
		if err != nil {
			log.Println("GetQueuedCompletionStatus err", err)
			return
		}
		log.Println("workThread new completion event,port", cpk.addr.Port)
		if ovPlus.isRead {
			log.Println("workThread is read")
			if bytesTrans == 0 {
				log.Println("workThread: conn closed~")
				continue
			}
			ovPlus.wsaBuf.Len = bytesTrans
			ovPlus.isRead = false
			err = windows.WSASend(cpk.sock.Fd, &ovPlus.wsaBuf, 1, nil, 0, ovPlus.ov, nil)
			if err != nil {
				log.Println("workThread: WSASend err", err)
				return
			}
			// 创建一个ioData用于读数据
			ovPlusForRead := &ioData{
				ov: &windows.Overlapped{},
				wsaBuf: windows.WSABuf{
					Len: BUF_SIZE,
					Buf: ovPlus.wsaBuf.Buf,
				},
				isRead: true,
			}
			err = windows.WSARecv(cpk.sock.Fd, &ovPlusForRead.wsaBuf, 1, nil, nil, ovPlusForRead.ov, nil)
			if err != nil {
				log.Println("workThread: WSARecv err", err)
				return
			}
		} else {
			log.Println("workThread: msg is sent!", ovPlus.wsaBuf.Len)
		}
	}
}

func (i *Iocp) AcceptLoop() {
	// 定义需要复用的变量
	var (
		b            = byte(0)
		wsaBuf       = windows.WSABuf{Len: BUF_SIZE, Buf: &b}
		recvBytesNum uint32
		flags        uint32
	)
	for {
		// 建立新连接
		clientFd, addr, err := windows.Accept(i.sock.Fd)
		if err != nil {
			log.Println("AcceptLoop err", err)
			return
		}
		// copy出来
		clientSock := *(addr.(*windows.SockaddrInet4))
		cpk := completionKey{
			sock: socket.NewSocket(clientFd),
			addr: &clientSock,
		}
		log.Println("AcceptLoop new conn, client port", cpk.addr.Port)

		// 监听这个新clientFd上的IO完成事件
		_, err = windows.CreateIoCompletionPort(clientFd, i.hComPort, uintptr(unsafe.Pointer(&cpk)), 0)
		if err != nil {
			log.Println("AcceptLoop CreateIoCompletionPort err", err)
			return
		}

		bit := byte(0)
		ovPlus := ioData{
			ov:     &windows.Overlapped{},
			isRead: true,
			wsaBuf: windows.WSABuf{Len: BUF_SIZE, Buf: &bit},
		} // 此参数不能填null，否则WSARecv将以阻塞方式工作
		// 非阻塞读取，内核读取完成后会调用GetQueuedCompletionStatus()通知到用户态
		err = windows.WSARecv(clientFd, &wsaBuf, 1, &recvBytesNum, &flags, ovPlus.ov, nil)
		if err != nil {
			//if windows.GetLastError() != windows.ERROR_IO_PENDING {
			//	log.Println("AcceptLoop WSARecv err", err)
			//	return
			//}
			log.Println("AcceptLoop WSARecv err", err)
			return
		}
	}
}
