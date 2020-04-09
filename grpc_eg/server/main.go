package main

import (
	"context"
	pb "github.com/chaseSpace/go-common-pkg-exmaples/grpc_eg/pb_test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"log"
	"net"
)

const (
	port = ":50051"
)

// serverSSS is used to implement helloworld.GreeterServer.
type serverSSS struct{}

// Search implements req.Search
func (s *serverSSS) Search(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Printf("ctx value a:%v\n", ctx.Value("a"))
	log.Printf("Received: query:%v header:%+v", in.GetQuery(), in.GetHeaders())
	return &pb.Response{ReqQuery: in.GetQuery(),
		X: &pb.ItemDetail{Name: "apple", Price: 110,
			Desc: "desc", Status: pb.ItemDetail_ACTIVE},
		Succ: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(testdata.Path("server1.pem"),
		testdata.Path("server1.key"))
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	svrOpt := grpc.Creds(creds)
	s := grpc.NewServer(svrOpt)
	pb.RegisterSearchSSSServer(s, &serverSSS{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
