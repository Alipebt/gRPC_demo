package main

import (
	"context"
	"fmt"
	"grpcDemo/pb"
	"log"
	"sync"
	"time"

	_ "github.com/apache/skywalking-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:5678", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	resp, err := client.SimpleRPC(context.Background(), &pb.HelloRequest{Name: "simpleRPC"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp.GetReply())

	req, err := client.BidrectionalStreamingRPC(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 开协程发送数据，发送10次
	go func() {
		for i := 0; i < 10; i++ {
			req.Send(&pb.HelloRequest{
				Name: "alipebt",
			})
			time.Sleep(time.Second)
		}
		wg.Done()
	}()
	// 开协程接收数据，接收10次
	go func() {
		for i := 0; i < 10; i++ {
			r, err := req.Recv()
			if err != nil {
				fmt.Println("出错了")
				return
			}
			fmt.Println(r.Reply)

		}
		wg.Done()
	}()
	wg.Wait()

	select {}
}
