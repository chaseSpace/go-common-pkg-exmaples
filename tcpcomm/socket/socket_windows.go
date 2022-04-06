package socket

import (
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"syscall"
)

type Socket struct {
	Fd int
}

func NewSocket(fd int) *Socket {
	return &Socket{
		Fd: fd,
	}
}

func (socket *Socket) Read(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(socket.Fd, bytes)
	if n < 0 {
		n = 0 // sometimes, n<0 is happening, it may cause caller panic
	}
	return n, err
}

func (socket *Socket) Write(bytes []byte) (int, error) {
	numBytesWritten, err := syscall.Write(socket.Fd, bytes)
	if err != nil {
		numBytesWritten = 0
	}
	return numBytesWritten, err
}

func (socket *Socket) Close() error {
	syscall.WSACleanup() // 释放为我初始化的dll资源
	return syscall.Close(socket.Fd)
}

// MakeWord 传入一个低位和高位字节数据参数，合并为一个双字节数据，作用同C++同名宏函数
func MakeWord(lowByte, highByte int8) uint16 {
	// 低8位，高8位
	low := uint16(lowByte)
	high := uint16(highByte) * 1 << 8
	return low + high
}

// LoByte 取一个双字节数据最低（最右边）字节的内容，作用同C++同名函数
func LoByte(x uint16) uint16 {
	x &= 0x000F // 高位清零就得到低位
	return x
}

// HiByte 取一个双字节数据最高（最左边）字节的内容，作用同C++同名函数
func HiByte(x uint16) uint16 {
	x &= 0xFF00 // 低位清零就得到高位
	return x
}

// WSA（Windows Sockets Asynchronous）文档
// https://docs.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-wsastartup

func Listen(ip string, port int) (sock *Socket, err error) {
	socket := &Socket{}

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
	wd := &syscall.WSAData{
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
	err = syscall.WSAStartup(uint32(socketVer), wd)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			syscall.WSACleanup() // 释放为我初始化的dll资源
		}
	}()
	// NOTE:这里还可以在检查wd返回的版本是否符合需求，不符合可退出
	// >>这里的wd.Version也是同上面MakeWord()方式计算出来的，所以要反解出主副版本，再比对
	if LoByte(wd.Version) != uint16(mainVer) || HiByte(wd.Version) != uint16(minorVer) {
		err = fmt.Errorf("could not find a usable version of Winsock.dll")
		return
	}
	h, err := windows.CreateIoCompletionPort(windows.InvalidHandle, windows.InvalidHandle, 0, 0)
	if err != nil {
		err = fmt.Errorf("could not find a usable version of Winsock.dll")
		return
	}

	// SOCK_STREAM 表示采用tcp协议
	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}
	err = syscall.SetNonblock(socketFileDescriptor, true)
	if err != nil {
		return nil, fmt.Errorf("failed to SetNonblock (%v)", err)
	}
	socket.Fd = socketFileDescriptor
	/*
		设置 SO_REUSEADDR 方便实现热重启
	*/
	err = syscall.SetsockoptInt(socket.Fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if err != nil {
		return nil, fmt.Errorf("failed set SO_REUSEADDR (%v)", err)
	}
	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	if err = syscall.Bind(socket.Fd, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	if err = syscall.Listen(socket.Fd, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil
}
