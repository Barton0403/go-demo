package main

import (
	"example.com/micro-go-book/ch7-rpc/grpc/string-service"
	"example.com/micro-go-book/ch7-rpc/pb"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
