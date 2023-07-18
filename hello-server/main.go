package main

import (
	"context"
	"fmt"
	"net"

	pb "main/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("service: SayHello")
	return &pb.HelloResponse{ResponseMsg: "hello " + req.RequestName}, nil
}

func main() {
	listen, _ := net.Listen("tcp", ":9090")

	grpcServer := grpc.NewServer()

	pb.RegisterSayHelloServer(grpcServer, &server{})

	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
