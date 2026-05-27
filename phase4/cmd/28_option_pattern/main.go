// Program 28_option_pattern 演示 Option 模式（带默认值的配置对象）。
//
// 与 Java Builder 模式对比：
//   - Java 常用链式 Builder（如 new ServerBuilder().setHost("0.0.0.0").setPort(8080).build()）；
//   - Go 的 Option 模式通常用函数选项（functional options）或结构体选项。
//   - Option 模式适合需要大量默认值、偶尔覆盖个别字段的场景。
package main

import "fmt"

// ServerConfig 服务器配置。
type ServerConfig struct {
	Host string
	Port int
	Mode string // debug / release
}

// Option 是修改 ServerConfig 的函数类型。
type Option func(*ServerConfig)

// WithHost 设置主机地址。
func WithHost(host string) Option {
	return func(c *ServerConfig) {
		c.Host = host
	}
}

// WithPort 设置端口。
func WithPort(port int) Option {
	return func(c *ServerConfig) {
		c.Port = port
	}
}

// WithMode 设置运行模式。
func WithMode(mode string) Option {
	return func(c *ServerConfig) {
		c.Mode = mode
	}
}

// NewServerConfig 创建配置，带有默认值。
func NewServerConfig(opts ...Option) *ServerConfig {
	cfg := &ServerConfig{
		Host: "127.0.0.1",
		Port: 8080,
		Mode: "debug",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func main() {
	fmt.Println("=== Option 模式演示 ===")

	// 使用全部默认值
	cfg1 := NewServerConfig()
	fmt.Printf("默认配置: %+v\n", cfg1)

	// 覆盖部分选项
	cfg2 := NewServerConfig(
		WithHost("0.0.0.0"),
		WithPort(9090),
		WithMode("release"),
	)
	fmt.Printf("自定义配置: %+v\n", cfg2)
}
