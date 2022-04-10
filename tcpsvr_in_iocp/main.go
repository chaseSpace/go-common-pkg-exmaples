//go:build windows
// +build windows

package main

import (
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_iocp/iocp"
	"os"
)

func main() {
	eventLoop, err := iocp.NewEventLoop("127.0.0.1", 8080)
	if err != nil {
		println("Failed to create event loop:", err)
		os.Exit(1)
	}
	defer eventLoop.Close()
	println("Server started. Waiting for incoming connections. ^C to exit.")

	eventLoop.AcceptLoop()
}
