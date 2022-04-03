//go:build linux
// +build linux

package socketmod

import (
	"fmt"
	"net"
	"strconv"
	"syscall"
)

type Socket struct {
	Fd int
}

func (socket Socket) Read(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	numBytesRead, err := syscall.Read(socket.Fd, bytes)
	if err != nil {
		numBytesRead = 0
	}

	return numBytesRead, err
}

func (socket Socket) Write(bytes []byte) (int, error) {
	numBytesWritten, err := syscall.Write(socket.Fd, bytes)
	if err != nil {
		numBytesWritten = 0
	}
	return numBytesWritten, err
}

func (socket *Socket) Close() error {
	return syscall.Close(socket.Fd)
}
func (socket *Socket) String() string {
	return strconv.Itoa(socket.Fd)
}

func Listen(ip string, port int) (*Socket, error) {
	socket := &Socket{}

	// SOCK_STREAM 表示采用tcp协议
	socketFileDescriptor, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}
	//err = syscall.SetNonblock(socketFileDescriptor, true)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to SetNonblock (%v)", err)
	//}
	socket.Fd = socketFileDescriptor
	/*
		设置 SO_REUSEADDR & SO_REUSEPORT 方便快速重启
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
