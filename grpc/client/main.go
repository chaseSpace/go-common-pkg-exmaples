package main

import (
	"context"
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/grpc/key"
	"github.com/chaseSpace/go-common-pkg-exmaples/grpc/pb_test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

var c *grpc.ClientConn

var client pb_test.SearchSSSClient

func init() {
	fmt.Println("client...")
	creds, err := credentials.NewClientTLSFromFile(key.Path("cert.crt"), "gxt.grpcsrv.auth")
	clientOpt := grpc.WithTransportCredentials(creds)

	// 指导：线上需要开启PermitWithoutStream，持续保持连接
	// Time建议设置10s，ping不要太密集，timeout设置3-5s就好，内网的话设置3s

	// --> 连接的最长生命周期由server端设置
	opt2 := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		// 连接无活动N段时间后去ping，最低10s，默认无限，也就是不会ping，断了也不知道
		Time: 10 * time.Second,
		// ping超时，建议5s(默认20s)，正常情况是马上收的到pong
		Timeout: 5 * time.Second,
		// true表示空闲时使用ping维持连接，false表示不会发送ping（上面两个参数被忽略）
		PermitWithoutStream: false,
	})

	c, err = grpc.Dial(address, clientOpt, opt2, grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println(c.GetState().String())
	client = pb_test.NewSearchSSSClient(c)
}

func call() {
	log.Println(111, c.GetState().String())
	// 不能通过ctx来传输kv
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)
	r, err := client.Search(context.WithValue(ctx, "a", 1),
		&pb_test.Request{Query: "q", Headers: map[string]string{"a": "B"}})
	if err != nil {
		// 当rpc error: code = Unavailable desc = transport is closing
		log.Printf("could not search: %v\n", err)
	} else {
		// parse timestamp
		log.Printf("searching: %+v", r.String())
	}
}

func main() {
	// Contact the server and print out its response.
	// 测试报告1：
	// server端宕机后，client在操作时总能立马变更连接状态(尝试重连，CONNECTING, TRANSIENT_FAILURE)
	// server恢复后，client也能马上恢复，无需重新拨号（除非client已经shutdown）

	// 报告2：当idle时间达到后，服务端会关闭连接
	// 此时client底层会自动无限重试（无论上层是否调用）
	defer c.Close()
	n := 0
	for n < 15 {
		call()
		//log.Println(111, c.GetState().String())
		time.Sleep(time.Second * 2)
		n++
	}
}
