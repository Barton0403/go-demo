package string_service

import (
	"context"
	"example.com/micro-go-book/ch7-rpc/pb"
)

type StringService struct {
	*pb.UnimplementedStringServiceServer
}

func (s *StringService) Concat(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	response := pb.StringResponse{Ret: "1"}
	return &response, nil
}

func (s *StringService) Diff(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	response := pb.StringResponse{Ret: "2"}
	return &response, nil
}
