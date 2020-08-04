package main

import (
	"context"
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/grpc/key"
	pb "github.com/chaseSpace/go-common-pkg-exmaples/grpc/pb_test"
	"github.com/golang/protobuf/ptypes"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"time"
)

const (
	// 建议server地址不要只填端口
	// 我遇到过server端只填端口，client必须填IP+Port才连得上的情况，反之时而可以时而不行
	addr = "localhost:50051"
)

// serverSSS is used to implement helloworld.GreeterServer.
type serverSSS struct{}

// Search implements req.Search
func (s *serverSSS) Search(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("Received: query:%v header:%+v", in.GetQuery(), in.GetHeaders())
	return &pb.Response{ReqQuery: in.GetQuery(),
		X: &pb.ItemDetail{Name: "apple", Price: 110,
			Desc: "desc", Status: pb.ItemDetail_ACTIVE},
		Succ: true,
		Time: ptypes.TimestampNow()}, nil
}

func (s *serverSSS) MockPanic(ctx context.Context, in *pb.MockPanicReq) (*pb.MockPanicRsp, error) {
	panic("here panic")
	return nil, nil
}

func main() {
	fmt.Println("server...")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(key.Path("cert.crt"),
		key.Path("rsa_private.key"))
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	svrOpt := grpc.Creds(creds)

	// 一般我们设置MaxConnectionIdle=5-10min足够了，其他无需更改
	// 服务器断开后，客户端会自动马上重连
	keepaliveopt := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute, // 空闲多久关闭连接，默认不关，0表示无限
		//MaxConnectionAge: 1*time.Second, // 连接最多存活多久，默认无限，不用设置
		//MaxConnectionAge+MaxConnectionAgeGrace这个时间过后，连接被强制关
		MaxConnectionAgeGrace: 1 * time.Second,
		Time:                  time.Hour,       // 空闲多久主动ping客户端，默认2hour
		Timeout:               1 * time.Second, // ping超时，默认20s
	})

	// 添加拦截器, 这里是有顺序的，recovery拦截器应在第一个，
	// 即接口报错应被立即捕捉到并返回err，后面的拦截器无需再执行（当然，这也不是绝对的，根据需要）
	interceptors := grpc_middleware.WithUnaryServerChain(
		RecoveryInterceptor,
		LoggingInterceptor)
	s := grpc.NewServer(
		svrOpt,
		keepaliveopt,
		interceptors,
	)
	pb.RegisterSearchSSSServer(s, &serverSSS{})

	defer func() {
		s.GracefulStop()
	}()
	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %v", err)
	}
}
