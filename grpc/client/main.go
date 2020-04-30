package main

import (
	"context"
	"github.com/chaseSpace/go-common-pkg-exmaples/grpc/pb_test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/alts/internal/conn"
	"google.golang.org/grpc/testdata"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

var c grpc.ClientConn

var client pb_test.SearchSSSClient

func init() {
	creds, err := credentials.NewClientTLSFromFile(testdata.Path("server1.pem"), "*.test.youtube.com")
	clientOpt := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial(address, clientOpt)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Println(conn.GetState())
	client = pb_test.NewSearchSSSClient(conn)
}

func main() {
	defer client()
	// Contact the server and print out its response.

	r, err := c.Search(context.WithValue(context.Background(), "a", 1),
		&pb_test.Request{Query: "q", Headers: map[string]string{"a": "B"}})
	if err != nil {
		log.Fatalf("could not search: %v", err)
	}
	// parse timestamp
	log.Print("Rsp time:", time.Unix(r.GetTime().Seconds, 0).String())
	log.Printf("searching: %+v", r.String())
}
