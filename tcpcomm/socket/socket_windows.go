package socket

import (
	"fmt"
	"golang.org/x/sys/windows"
	"net"
	"syscall"
)

type Socket struct {
	Fd windows.Handle
}

func NewSocket(fd windows.Handle) *Socket {
	return &Socket{
		Fd: fd,
	}
}

func (socket *Socket) Read(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	n, err := windows.Read(socket.Fd, bytes)
	if n < 0 {
		n = 0 // sometimes, n<0 is happening, it may cause caller panic
	}
	return n, err
}

func (socket *Socket) Write(bytes []byte) (int, error) {
	numBytesWritten, err := windows.Write(socket.Fd, bytes)
	if err != nil {
		numBytesWritten = 0
	}
	return numBytesWritten, err
}

func (socket *Socket) Close() error {
	return windows.Close(socket.Fd)
}

// WSA（Windows Sockets Asynchronous）文档
// https://docs.microsoft.com/en-us/windows/win32/api/winsock/nf-winsock-wsastartup

func Listen(ip string, port int) (sock *Socket, err error) {
	socket := &Socket{}

	sockHandle, err := windows.WSASocket(windows.AF_INET, windows.SOCK_STREAM, 0, nil,
		0, windows.WSA_FLAG_OVERLAPPED)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket (%v)", err)
	}
	socket.Fd = sockHandle
	defer func() {
		if err != nil {
			_ = socket.Close()
		}
	}()
	socketAddress := &windows.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	if err = windows.Bind(socket.Fd, socketAddress); err != nil {
		return nil, fmt.Errorf("failed to bind socket (%v)", err)
	}

	if err = windows.Listen(socket.Fd, syscall.SOMAXCONN); err != nil {
		return nil, fmt.Errorf("failed to listen on socket (%v)", err)
	}

	return socket, nil
}
