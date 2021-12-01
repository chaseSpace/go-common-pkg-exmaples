package main

import (
	"log"
	"net"
)

func main() {
	_, err := net.Listen("tcp", "192.168.1.1:8080")
	if err != nil {
		log.Fatal(111, err)
	}
	_, err = net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 1),
		Port: 8080,
		Zone: "",
	})
	if err != nil {
		log.Fatal(222, err)
	}
}

type S struct {
	X string
}
