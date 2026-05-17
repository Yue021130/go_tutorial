// Package gmp 包含 GMP 调度模型相关的小工具函数。
package gmp

import (
	"fmt"
	"runtime"
)

// PrintGMPInfo 打印与 GMP 相关的运行时信息。
func PrintGMPInfo() {
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU:     %d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine: %d\n", runtime.NumGoroutine())
}

// GOMAXPROCS 设置并返回 P 的数量。
func GOMAXPROCS(n int) int {
	return runtime.GOMAXPROCS(n)
}
