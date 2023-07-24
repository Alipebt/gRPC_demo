package main

import (
	"context"
	"fmt"
	"grpcDemo/pb"
	"log"
	"net"
	"sync"
	"time"

	_ "github.com/apache/skywalking-go"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SimpleRPC(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("client call simpleRPC...")
	log.Println(in)
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func (s *server) AllStream(request pb.Greeter_BidrectionalStreamingRPCServer) error {
	// 开启两个协程，一个接收数据，一个发送数据
	var wg = sync.WaitGroup{}
	wg.Add(2)
	go func() { // 发送数据
		for {
			err := request.Send(&pb.HelloResponse{
				Reply: "Send Hello",
			})
			if err != nil {
				fmt.Println(err)
				break
			}
			time.Sleep(time.Second) // 睡1s发送一次
		}
		wg.Done()

	}()
	go func() { // 不停接收数据
		for {
			res, err := request.Recv()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(res.Name)
		}
		wg.Done()

	}()

	wg.Wait()
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":5678")
	if err != nil {
		log.Fatal(err)
		return
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	log.Println("gRPC server starts running...")

	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
		return
	}
}
