// Program 26_grpc_server 演示一个 gRPC 服务端。
//
// 与 REST 对比：
//   - REST 基于 HTTP/1.1 或 HTTP/2 + JSON，通用性好，调试方便；
//     gRPC 基于 HTTP/2 + Protocol Buffers，性能更高，强类型契约。
//   - REST 适合面向浏览器/第三方开放 API；gRPC 适合微服务内部通信。
//
// 与 Java Dubbo 对比：
//   - Dubbo 基于 TCP + Hessian/Protobuf，依赖注册中心（ZooKeeper/Nacos）；
//   - gRPC 基于 HTTP/2 + Protobuf，可配合 etcd/Consul/Nacos 做服务发现。
//
// 本示例使用 grpc-go 官方示例包中的 helloworld 服务，避免本地安装 protoc。
// 生产环境中，应自己编写 .proto 文件并通过 protoc 生成 pb.go。
package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// server 实现 helloworld.GreeterServer 接口。
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello 实现 SayHello RPC 方法。
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}

func main() {
	addr := "127.0.0.1:19004"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen failed: %v\n", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})

	fmt.Printf("gRPC Server 启动: %s\n", addr)
	fmt.Println("运行 27_grpc_client 进行测试")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("serve failed: %v\n", err)
	}
}
