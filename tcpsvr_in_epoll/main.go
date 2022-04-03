//go:build linux
// +build linux

package main

import (
	"bufio"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_epoll/epollmod"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_epoll/socketmod"
	"log"
	"os"
	"strings"
)

func main() {
	eventLoop, err := epollmod.NewEventLoop("127.0.0.1", 8080)
	if err != nil {
		log.Println("Failed to create event loop:", err)
		os.Exit(1)
	}
	defer eventLoop.Close()
	log.Println("Server started. Waiting for incoming connections. ^C to exit.")
	eventLoop.Handle(func(s *socketmod.Socket) {
		reader := bufio.NewReader(s)
		for {
			//b, err := reader.ReadByte()
			//log.Println("incoming data...", b)
			line, err := reader.ReadString('\n')
			if err != nil || strings.TrimSpace(line) == "" {
				break
			}
			log.Println("incoming data...", strings.TrimRight(line, "\n"))
			s.Write([]byte(line))
		}
		s.Close()
	})
}
