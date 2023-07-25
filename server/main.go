package main

import (
	"context"
	"grpcDemo/pb"
	"log"
	"net"

	_ "github.com/apache/skywalking-go"
	"google.golang.org/grpc"
)

type Echo struct {
	pb.UnimplementedEchoServer
}

//  ServerStreamingEcho 客户端发送一个请求 服务端以流的形式循环发送多个响应
/*
1. 获取客户端请求参数
2. 处理完成后返回过个响应
3. 最后返回nil表示已经完成响应
*/
func (e *Echo) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	log.Printf("Recved %v", req.GetMessage())
	// 具体返回多少个response根据业务逻辑调整
	for i := 0; i < 2; i++ {
		// 通过 send 方法不断推送数据
		err := stream.Send(&pb.EchoResponse{Message: req.GetMessage()})
		if err != nil {
			log.Fatalf("Send error:%v", err)
			return err
		}
	}
	// 返回nil表示已经完成响应
	return nil
}

// UnaryEcho 一个普通的UnaryAPI
func (e *Echo) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("Recved: %v", req.GetMessage())
	resp := &pb.EchoResponse{Message: req.GetMessage()}
	return resp, nil
}

func main() {
	listen, err := net.Listen("tcp", ":5678")
	if err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterEchoServer(s, &Echo{})
	log.Println("gRPC server starts running...")

	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
}
