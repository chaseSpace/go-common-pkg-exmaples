//go:build linux
// +build linux

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_epoll/epollmod"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpsvr_in_epoll/socketmod"
	"io"
	"log"
	"os"
	"strings"
	"syscall"
	"time"
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
		log.Println("eventLoop.Handle start ======")
		// 下面把所有收到的数据返回去，模拟的HTTP response
		b := bytes.Buffer{}
		for {
			//b, err := reader.ReadByte()
			//log.Println("incoming data...", b)
			line, err := reader.ReadString('\n')
			if err == nil {
				b.WriteString(line)
				log.Println("Handle incoming data...", strings.TrimRight(line, "\n"))
				continue
			}
			if err == syscall.EAGAIN { // 当前缓冲区已无数据可读
				log.Println("Handle data EOF")
				break
			}
			if err == syscall.EINTR || err == syscall.EBADF { // 可忽略的错误
				continue
			}
			if err == io.ErrNoProgress { // 读完了
				break
			}
			// 其他无法处理的错误
			_ = s.Close()
			log.Println("Handle other err:", err)
			return
		}
		if b.Len() == 0 {
			// 对方已关闭conn
			log.Println("Handle peer closed socket")
			s.Close()
		}
		body := fmt.Sprintf(`<html>
								  <head>
									<title>Epoll Response</title>
								  </head>
								  <body>%s</body>
								</html>`, b.Bytes())
		fmt.Fprintf(s, `HTTP/1.1 200 OK
        Content-Type: text/html;charset=UTF-8
        Content-Length: %d
        Date: %s
		%s`, len(body), time.Now().Format(time.RFC1123), body)
		s.Close()
	})
}
