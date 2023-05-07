package main

import (
	"context"
	"fmt"
	"net"

	pb "grpc_demo/server/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSendMsgServer
}

func (s *server) SendMsg(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	fmt.Printf("received: %v\n", req.RequestMsg)
	return &pb.Response{ResponseMsg: "response: " + req.RequestMsg + "\n"}, nil
}

func main() {
	//监听端口
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Printf("listen port error: %v\n", err)
		return
	}
	//新建grpc的服务
	grpcServer := grpc.NewServer()
	//注册服务
	pb.RegisterSendMsgServer(grpcServer, &server{})
	//启动
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("start grpcServer error: %v\n", err)
		return
	}
}
