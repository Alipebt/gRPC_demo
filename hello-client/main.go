package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "main/proto"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSayHelloClient(conn)

	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Alipebt"})

	fmt.Println(resp.GetResponseMsg())

}
