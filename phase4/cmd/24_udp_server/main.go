// Program 24_udp_server 演示用 Go 标准库实现一个 UDP Echo 服务器。
//
// 与 TCP 区别：
//   - TCP 是面向连接、可靠传输；UDP 是无连接、尽力交付。
//   - TCP 用 Accept + 每个连接一个 goroutine；UDP 直接 ReadFrom/WriteTo 处理数据报。
//   - UDP 适合实时性要求高、可容忍丢包的场景（如 DNS、音视频传输、日志上报）。
package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	addr := "127.0.0.1:19002"
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Printf("resolve failed: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Printf("listen failed: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("UDP Echo Server 启动: %s\n", addr)
	fmt.Println("可用 nc -u 127.0.0.1 19002 测试")

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read failed: %v\n", err)
			continue
		}

		msg := strings.TrimSpace(string(buf[:n]))
		fmt.Printf("收到来自 %s: %s\n", clientAddr, msg)

		// 直接回显
		reply := fmt.Sprintf("UDP Echo: %s\n", msg)
		_, _ = conn.WriteToUDP([]byte(reply), clientAddr)
	}
}
