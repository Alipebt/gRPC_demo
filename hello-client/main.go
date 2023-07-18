package main

import (
	"context"
	"fmt"
	pb "grpc_demo/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, _ := grpc.Dial("127.0.0.1:9090")
	defer conn.Close()

	client := pb.NewSayHelloClient(conn)
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "alipebt "})
	fmt.Println(resp.GetRespanseMsg())
}
