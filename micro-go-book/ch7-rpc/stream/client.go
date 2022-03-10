package main

import (
	"context"
	"example.com/micro-go-book/ch7-rpc/stream_pb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()

	bookClient := stream_pb.NewStringServiceClient(conn)
	stringReq := &stream_pb.StringRequest{A: "A", B: "B"}
	stream, _ := bookClient.LotsOfServerStream(context.Background(), stringReq)
	for {
		item, stream_err := stream.Recv()
		if stream_err == io.EOF {
			break
		}
		if stream_err != nil {
			log.Printf("failed to recv: %v", stream_err)
		}
		fmt.Printf("StringService Concat: %s concat %s = %s\n", stringReq.A, stringReq.B, item.GetRet())
	}
}

func sendClientStreamRequest(client stream_pb.StringServiceClient) {

}
