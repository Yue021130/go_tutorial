// Program 25_http2_server 演示 HTTP/2 服务器（h2c，即无 TLS 的 HTTP/2）。
//
// 与 HTTP/1.1 对比：
//   - HTTP/2 支持多路复用（单个 TCP 连接上并行传输多个请求/响应）。
//   - HTTP/2 支持头部压缩（HPACK）。
//   - HTTP/2 支持服务器推送（Server Push）。
//
// 与 Java 对比：
//   - Java 需要 Tomcat 9+ / Jetty 9.3+ / Undertow 等容器显式开启 HTTP/2，
//     且通常需要 TLS + ALPN。
//   - Go 标准库 net/http 在 TLS 下自动协商 HTTP/2；非 TLS 场景可通过
//     golang.org/x/net/http2/h2c 启用 h2c。
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "protocol: %s\nmethod: %s\npath: %s\n",
			r.Proto, r.Method, r.URL.Path)
	})

	// h2cHandler 允许客户端通过明文 HTTP/2 直接连接
	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    ":19003",
		Handler: h2c.NewHandler(mux, h2s),
	}

	fmt.Println("HTTP/2 (h2c) Server 启动: http://127.0.0.1:19003")
	fmt.Println("测试命令：curl --http2 -v http://127.0.0.1:19003/")
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server failed: %v\n", err)
	}
}
