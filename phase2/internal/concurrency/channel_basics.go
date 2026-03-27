package concurrency

import (
	"fmt"
	"sync"
	"time"
)

// ==================== buffered vs unbuffered channel ====================
//
// unbuffered channel: make(chan int)
//   - 发送和接收必须同时准备好，是同步的
//   - 发送方会阻塞直到有接收方
//   - 用于 goroutine 之间的直接交接
//
// buffered channel: make(chan int, n)
//   - 发送方只在缓冲区满时阻塞
//   - 接收方只在缓冲区空时阻塞
//   - 用于解耦生产者和消费者速率

func DemonstrateUnbuffered() {
	ch := make(chan string)
	go func() {
		fmt.Println("goroutine: 准备发送")
		ch <- "hello" // 会阻塞，直到主 goroutine 接收
		fmt.Println("goroutine: 发送完成")
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("main: 准备接收")
	msg := <-ch
	fmt.Println("main: 收到", msg)
}

func DemonstrateBuffered() {
	ch := make(chan string, 2)
	ch <- "one" // 不阻塞，缓冲区未满
	ch <- "two" // 不阻塞
	fmt.Println("main: 已发送两个值，未被接收")

	fmt.Println("main: 收到", <-ch)
	fmt.Println("main: 收到", <-ch)
}

// ==================== channel 关闭与 range ====================
//
// 关闭 channel 后：
//   - 继续发送会 panic
//   - 继续接收会返回零值和 false（comma-ok 模式）
//   - range 会在 channel 关闭后自动结束

func GenerateWithClose(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			out <- i
		}
	}()
	return out
}

// ==================== sync.WaitGroup ====================
//
// WaitGroup 用于等待一组 goroutine 完成。
// 与 Java 的 CountDownLatch 类似。

func WaitGroupExample() {
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("worker %d starting\n", id)
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("worker %d done\n", id)
		}(i)
	}
	wg.Wait()
	fmt.Println("all workers done")
}

// ==================== goroutine 泄漏示例 ====================
//
// 泄漏原因：goroutine 永远阻塞在某个 channel 操作上，无法退出。
// 解决方案：使用 context 取消或关闭 channel。

func LeakyGoroutine() {
	ch := make(chan int)
	go func() {
		// 这个 goroutine 会永远阻塞在这里，因为没有人发送数据
		v := <-ch
		fmt.Println("never printed:", v)
	}()
	// main 不等待 goroutine 结束，导致泄漏
}

// SafeGoroutine 展示了如何用 context 避免泄漏
// （具体实现放在 cancellation 包中）

// ==================== 死锁模式 ====================
//
// 死锁：所有 goroutine 都在等待，没有人在推进。
// Go 运行时会在检测到死锁时 panic。
//
// 常见死锁：
// 1. 无缓冲 channel 上只有发送没有接收
// 2. channel 已满，发送方阻塞，接收方也阻塞
// 3. 锁的循环获取

// DeadlockExample 故意制造一个死锁（仅用于演示，不要运行）
// func DeadlockExample() {
//     ch := make(chan int)
//     ch <- 1 // 没有接收方，死锁
// }
