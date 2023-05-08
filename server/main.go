package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	pb "grpc_demo/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type server struct {
	pb.UnimplementedSendMsgServer
}

func validateToken(ctx context.Context) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	tokens, ok := md["key"]
	if !ok || len(tokens) == 0 {
		return errors.New("no token is provided")
	}

	if tokens[0] != "value" {
		return errors.New("invalid token")
	}

	return nil
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	err := validateToken(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return handler(ctx, req)

}

func (s *server) SendMsg(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	fmt.Printf("received: %v\n", req.RequestMsg)

	return &pb.Response{ResponseMsg: "response: " + req.RequestMsg + "\n"}, nil

}

func main() {
	// SSL/TSL安全认证
	creds, _ := credentials.NewServerTLSFromFile(
		"/home/alipebt/code/grpc_demo/ssl/own.pem",
		"/home/alipebt/code/grpc_demo/ssl/own.key")

	// 监听端口
	listen, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Printf("listen port error: %v\n", err)
		return
	}

	// 新建grpc的服务
	grpcServer := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor))

	// 注册服务
	pb.RegisterSendMsgServer(grpcServer, &server{})

	// 启动
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("start grpcServer error: %v\n", err)
		return
	}
}
