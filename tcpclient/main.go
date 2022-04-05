package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

/*
这个tcp client是为此项目下的
	-	tcpsvr_in_epoll_LT
	-	tcpsvr_in_epoll_ET
编写的
*/
const packDelimiter = '\n'

func genContent(length int) string {
	buf := bytes.NewBufferString("start_")
	for i := 0; i < length; i++ {
		buf.WriteString("0123456789")
	}
	buf.WriteString("_end")
	return buf.String()
}

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	println("conn ok...")

	serverDown := make(chan bool)
	go func() {
		buf := bufio.NewReader(conn)
		for {
			resp, err := buf.ReadString(packDelimiter)
			if err == io.EOF {
				println("conn closed by server")
				serverDown <- true
				return
			}
			if err != nil {
				println("recv err", err)
				return
			}
			resp = strings.TrimRight(resp, string(packDelimiter))
			println(resp)
		}
	}()
	// 作为用于简单测试连通性的tcp client，此处仅发送几次与server端约定好的data pack，观察回包即可
	for i := 0; i < 1; i++ {
		// --	这里一次发送2个pack，是为了测试server代码是否能够正常解析pack，#相关逻辑在server的conn.ReadPack()
		pack1 := fmt.Sprintf("client sent %s%s", genContent(10), string(packDelimiter))
		pack2 := fmt.Sprintf("client sent %s%s", genContent(10), string(packDelimiter))

		n, err := fmt.Fprintf(conn, pack1+pack2)
		fmt.Printf("write %d err:%v\n", n, err)
		time.Sleep(time.Second)
	}
	println("nothing to do...")
	defer conn.Close()

	select {
	case <-serverDown:
	}
}
