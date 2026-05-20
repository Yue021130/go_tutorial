// Program 22_escape_analysis 演示 Go 的内存逃逸分析、栈扩容与 sync.Pool。
//
// 核心概念：
//   - 逃逸分析（Escape Analysis）：编译器决定变量分配在栈上还是堆上的过程。
//   - 栈分配：函数返回后自动回收，速度极快；适合生命周期随函数结束的对象。
//   - 堆分配：可被函数外部引用，由 GC 回收；适合生命周期不确定或共享对象。
//   - 栈扩容：goroutine 栈初始 2KB，按需增长，最大可达 1GB（64 位系统）。
//   - sync.Pool：对象池，用于复用临时对象，减少 GC 压力。
//
// 栈扩容示意：
//
//   初始栈（2KB）        需要更多空间           扩容后栈
//   +----------+        +----------+          +----------+
//   |          |        |          |          |          |
//   |  小帧    |  -->   |  溢出！  |   -->    |  更大帧  |
//   |          |        |          |          |          |
//   +----------+        +----------+          +----------+
//
// 与 Java 对比：
//   - Java 对象几乎都在堆上分配（除标量替换等 JIT 优化）；Go 优先栈分配，
//     通过逃逸分析把尽可能多的对象留在栈上，减少 GC 负担。
//   - Java 线程栈大小固定（-Xss）；Go goroutine 栈动态伸缩，按需分配。
//   - Java 有对象池（Apache Commons Pool）；Go 标准库提供 sync.Pool，
//     但更轻量，且池中对象可能在 GC 时被回收，不能当作持久缓存。
package main

import (
	"fmt"
	"sync"
	"time"

	"go-tutorial/phase4/internal/escape"
)

func main() {
	fmt.Println("=== 内存逃逸、栈扩容与 sync.Pool 演示 ===")

	escape.PrintEscapeHint()
	fmt.Println()

	// 栈分配：x 通常不逃逸
	v := escape.StackValue()
	fmt.Printf("StackValue 返回值: %d\n", v)

	// 堆分配：返回指针导致 x 逃逸
	p := escape.HeapValue()
	fmt.Printf("HeapValue 返回值: %d\n", *p)

	// 切片逃逸到底层数组会分配在堆上
	s := escape.EscapeSlice(10)
	fmt.Printf("EscapeSlice: %v\n", s)

	// 使用外部切片避免额外分配
	dst := make([]int, 10)
	escape.NoEscapeSlice(dst)
	fmt.Printf("NoEscapeSlice: %v\n", dst)

	// 闭包逃逸
	f := escape.ClosureEscape()
	fmt.Printf("ClosureEscape: %d\n", f())

	// sync.Pool 演示：复用临时对象
	fmt.Println("\n=== sync.Pool 对象复用演示 ===")
	pool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("  Pool.New: 创建新对象")
			return make([]byte, 1024)
		},
	}

	// 从池中取对象
	obj1 := pool.Get().([]byte)
	obj1[0] = 'A'
	fmt.Printf("从池中取出对象，obj1[0]=%c\n", obj1[0])

	// 放回池中复用
	pool.Put(obj1)

	// 再次取出，可能是同一个对象
	obj2 := pool.Get().([]byte)
	fmt.Printf("再次取出对象，obj2[0]=%c（可能是之前放回的）\n", obj2[0])

	// 连续取多次，观察 New 调用次数
	for i := 0; i < 5; i++ {
		b := pool.Get().([]byte)
		_ = b
		pool.Put(b)
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("\n=== 演示结束 ===")
	fmt.Println("提示：运行 'go run -gcflags=-m ./cmd/22_escape_analysis' 查看逃逸分析输出")
}
