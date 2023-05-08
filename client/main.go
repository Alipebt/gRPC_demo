package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "grpc_demo/proto"
)

type ClientTokenAuth struct {
}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context,uri ...string) (map[string]string, error) {

	return map[string]string{
		"key": "value",
	}, nil

}

func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return true
}

func main() {
	// SSL/TSL安全认证
	creds, _ := credentials.NewClientTLSFromFile(
		"/home/alipebt/code/grpc_demo/ssl/own.pem",
		"*.grpcdemo.com")

	// 连接到server
	var opts []grpc.DialOption
	opts = append(opts,grpc.WithTransportCredentials(creds))
	opts = append(opts,grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	conn, err := grpc.Dial("127.0.0.1:9999", opts...)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer conn.Close()

	// 创建client
	client := pb.NewSendMsgClient(conn)

	// 调用远程方法
	resp, err := client.SendMsg(context.Background(),&pb.Request{RequestMsg: "massages"})
	fmt.Printf(resp.GetResponseMsg())

	return
}
