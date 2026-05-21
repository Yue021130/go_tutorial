// Program 23_tcp_server 演示用 Go 标准库实现一个简易 TCP Echo 服务器。
//
// 与 Java 对比：
//   - Java 通常用 ServerSocket + ThreadPool 处理并发连接；
//     Go 用 net.Listen + goroutine per connection，代码更简洁。
//   - Java 的 IO 有阻塞 BIO、NIO、AIO 之分；Go 的 net 包基于 epoll/kqueue/IOCP
//     的 runtime 网络轮询器（netpoller），对开发者屏蔽了底层差异，
//     统一使用阻塞式 API + goroutine 模型。
package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	addr := "127.0.0.1:19001"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("listen failed: %v\n", err)
		return
	}
	defer listener.Close()
	fmt.Printf("TCP Echo Server 启动: %s\n", addr)
	fmt.Println("可用 telnet/nc 测试：nc 127.0.0.1 19001")

	// 接受连接循环
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept failed: %v\n", err)
			continue
		}

		// 每个连接一个 goroutine，无需线程池
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("客户端连接: %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	for {
		// 设置读超时，避免连接一直挂起
		_ = conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("客户端断开: %s, err=%v\n", conn.RemoteAddr(), err)
			return
		}

		msg := strings.TrimSpace(line)
		if msg == "quit" {
			fmt.Fprintf(conn, "Bye!\n")
			return
		}

		// Echo 回显
		fmt.Fprintf(conn, "Echo: %s\n", msg)
	}
}
