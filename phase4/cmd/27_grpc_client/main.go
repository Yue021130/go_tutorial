// Program 27_grpc_client 演示一个 gRPC 客户端。
package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	addr := "127.0.0.1:19004"

	// 创建连接（生产环境应配置 TLS）
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("connect failed: %v\n", err)
		return
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Go Learner"})
	if err != nil {
		fmt.Printf("say hello failed: %v\n", err)
		return
	}
	fmt.Printf("gRPC 响应: %s\n", r.GetMessage())
}
