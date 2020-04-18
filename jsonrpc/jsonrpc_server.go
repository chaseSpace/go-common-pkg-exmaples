package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

/*
rpc 是一个远程调用通信协议。
处理逻辑位于远端，客户端只需要知道对应函数，参数，返回值格式即可，不关心执行过程，只关心结果。
根据实现不同，可基于 tcp/udp/http传输

golang中，典型的rpc库有：
1. 标准库 net/rpc  [基于tcp/http 传输方式] !!!_+_ 不支持跨语言调用
2. net/rpc/jsonrpc 支持跨语言<<<<<<<<<<<<<<<<<<<
2. gPRC 快，支持跨语言
*/

// 在server端

type Listener int

type Reply struct {
	Data string
}

func (l *Listener) GetLine(line []byte, reply *Reply) error {
	rv := string(line)
	fmt.Printf("Receive: %v\n", rv)
	*reply = Reply{rv}
	return nil
}

func (l *Listener) Sayhi(data []byte, reply *Reply) error {
	rv := string(data)
	fmt.Printf("Receive: %v\n", rv)
	if rv == "lei" || rv == "lucy" {
		(*reply).Data = "hi, " + rv
		return nil
	}
	return errors.New(fmt.Sprintf("data must be `lei` or `lucy`, get{%s}", rv))
}

func main() {
	addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:12345")
	if err != nil {
		log.Fatal(err)
	}

	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		log.Fatal(err)
	}

	listener := new(Listener)
	err = rpc.Register(listener)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := inbound.Accept()
		if err != nil {
			log.Fatal(err)
		}
		jsonrpc.ServeConn(conn)
	}
}
