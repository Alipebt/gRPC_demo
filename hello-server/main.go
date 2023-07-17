package main

import (
	"context"
	pb "grpc_demo/proto"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{RespanseMsg: "hello" + req.RequestName}, nil
}

func main() {

	listen, _ := net.Listen("tcp", ":9090")
	grpcServer := grpc.NewServer()
	pb.RegisterSayHelloServer(grpcServer, &server{})
	err := grpcServer.Serve(listen)
	if err != nil {
		return
	}
	return
}
