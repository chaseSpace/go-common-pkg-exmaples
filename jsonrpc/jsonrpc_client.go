package main

import (
	"bufio"
	"log"
	"net/rpc/jsonrpc"
	"os"
)

type Reply struct {
	Data string
}

func main() {
	client, err := jsonrpc.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}

	input := bufio.NewReader(os.Stdin)
	for {
		// 接收终端输入(阻塞)
		line, _, err := input.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		var reply Reply
		err = client.Call("Listener.Sayhi", line, &reply)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Reply: %v, Data: %v", reply, reply.Data)
	}
}
