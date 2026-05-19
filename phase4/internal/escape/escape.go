// Package escape 用于演示 Go 的内存逃逸分析相关概念。
package escape

import "fmt"

// StackValue 返回一个在栈上分配的整数值（通常不逃逸）。
func StackValue() int {
	x := 42
	return x
}

// HeapValue 返回一个指针，导致 x 逃逸到堆上。
func HeapValue() *int {
	x := 42
	return &x
}

// EscapeSlice 返回一个切片，底层数组会逃逸到堆上。
func EscapeSlice(n int) []int {
	data := make([]int, n)
	for i := range data {
		data[i] = i
	}
	return data
}

// NoEscapeSlice 接收外部切片，避免内部分配逃逸。
func NoEscapeSlice(dst []int) {
	for i := range dst {
		dst[i] = i * 2
	}
}

// ClosureEscape 演示闭包捕获变量导致的逃逸。
func ClosureEscape() func() int {
	x := 100
	return func() int {
		return x
	}
}

// PrintEscapeHint 打印逃逸分析观察提示。
func PrintEscapeHint() {
	fmt.Println("使用 'go build -gcflags=-m' 可查看编译器逃逸分析结果")
	fmt.Println("示例：go build -gcflags=-m ./cmd/22_escape_analysis")
}
