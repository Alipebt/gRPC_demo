package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc_demo/server/proto"
)

func main() {
	//连接到server
	conn, err := grpc.Dial("127.0.0.1:9999",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	//创建client
	client := pb.NewSendMsgClient(conn)
	//调用远程方法
	resp, err := client.SendMsg(context.Background(), &pb.Request{RequestMsg: "massages"})
	fmt.Printf(resp.GetResponseMsg())

	return
}
