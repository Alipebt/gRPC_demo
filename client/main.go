package main

import (
	"context"
	"fmt"
	"grpcDemo/pb"
	"io"
	"log"

	_ "github.com/apache/skywalking-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func serverStream(client pb.EchoClient) {
	// 2.调用获取stream
	stream, err := client.ServerStreamingEcho(context.Background(), &pb.EchoRequest{Message: "Hello World"})
	if err != nil {
		log.Fatalf("could not echo: %v", err)
	}

	// 3. for循环获取服务端推送的消息
	for {
		// 通过 Recv() 不断获取服务端send()推送的消息
		resp, err := stream.Recv()
		// 4. err==io.EOF则表示服务端关闭stream了 退出
		if err == io.EOF {
			log.Println("server closed")
			break
		}
		if err != nil {
			log.Printf("Recv error:%v", err)
			continue
		}
		log.Printf("Recv data:%v", resp.GetMessage())
	}
}

func unary(client pb.EchoClient) {
	resp, err := client.UnaryEcho(context.Background(), &pb.EchoRequest{Message: "hello world"})
	if err != nil {
		log.Printf("send error:%v\n", err)
	}
	fmt.Printf("Recved:%v \n", resp.GetMessage())
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:5678", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	client := pb.NewEchoClient(conn)

	serverStream(client)
	// unary(client)

	select {}
}
