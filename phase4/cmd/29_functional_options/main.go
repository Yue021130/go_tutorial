// Program 29_functional_options 演示 Functional Options 模式。
//
// 与 Java 对比：
//   - Java 的构造器或 Builder 模式需要显式定义大量方法；
//   - Go 的 Functional Options 模式通过变长函数参数实现，API 简洁、扩展性好，
//     被广泛应用于 Go 标准库和知名项目（如 grpc.NewServer(opts...)、zap.New(core, opts...)）。
//
// 推荐实践（来自 Dave Cheney 的 Functional options for friendly APIs）：
//   - 选项函数返回一个未导出类型，避免外部直接构造。
//   - 提供合理的默认值。
//   - 选项函数名以 With 或 Set 开头。
package main

import (
	"fmt"
	"time"
)

// Server 是要配置的服务器对象。
type Server struct {
	host     string
	port     int
	timeout  time.Duration
	maxConns int
}

// Option 是未导出类型，保证 API 可扩展且不被外部破坏。
type Option func(*Server)

// WithHost 设置监听地址。
func WithHost(host string) Option {
	return func(s *Server) { s.host = host }
}

// WithPort 设置端口。
func WithPort(port int) Option {
	return func(s *Server) { s.port = port }
}

// WithTimeout 设置超时。
func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

// WithMaxConns 设置最大连接数。
func WithMaxConns(n int) Option {
	return func(s *Server) { s.maxConns = n }
}

// NewServer 创建 Server，应用 Functional Options。
func NewServer(opts ...Option) *Server {
	s := &Server{
		host:     "127.0.0.1",
		port:     8080,
		timeout:  30 * time.Second,
		maxConns: 100,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) String() string {
	return fmt.Sprintf("Server{host=%s, port=%d, timeout=%v, maxConns=%d}",
		s.host, s.port, s.timeout, s.maxConns)
}

func main() {
	fmt.Println("=== Functional Options 模式演示 ===")

	s1 := NewServer()
	fmt.Println(s1)

	s2 := NewServer(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithTimeout(60*time.Second),
		WithMaxConns(1000),
	)
	fmt.Println(s2)
}
