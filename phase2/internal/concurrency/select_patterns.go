package concurrency

import (
	"fmt"
	"time"
)

// ==================== select 随机选择 ====================
//
// 当多个 case 同时就绪时，select 会随机选择一个执行。
// 这是为了避免某个 channel 一直饥饿。

func SelectRandomDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() { ch1 <- "from ch1" }()
	go func() { ch2 <- "from ch2" }()

	// 两个 channel 几乎同时就绪，select 随机选择
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("received", msg)
		case msg := <-ch2:
			fmt.Println("received", msg)
		}
	}
}

// ==================== select with default：非阻塞操作 ====================
//
// default 分支会在没有 case 可执行时立即执行，实现非阻塞 send/receive。

func NonBlockingReceive(ch chan string) {
	select {
	case msg := <-ch:
		fmt.Println("received:", msg)
	default:
		fmt.Println("no message available")
	}
}

func NonBlockingSend(ch chan string, msg string) {
	select {
	case ch <- msg:
		fmt.Println("sent:", msg)
	default:
		fmt.Println("channel full, dropped:", msg)
	}
}

// ==================== ticker 与 timer ====================
//
// time.Ticker 周期触发，time.Timer 单次触发。
// 记得 Stop ticker，避免 goroutine 泄漏。

func TickerDemo(duration time.Duration) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	done := time.After(duration)
	for {
		select {
		case t := <-ticker.C:
			fmt.Println("tick at", t.Format("15:04:05.000"))
		case <-done:
			fmt.Println("ticker demo done")
			return
		}
	}
}

// ==================== 优雅关闭（graceful shutdown）====================
//
// 模式：一个 done channel 通知所有 goroutine 退出，
// 主 goroutine 等待所有工作 goroutine 完成后再退出。

func GracefulShutdown(jobs []int) {
	done := make(chan struct{})
	result := make(chan int, len(jobs))

	for _, job := range jobs {
		go func(j int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("recovered from panic:", r)
				}
			}()

			select {
			case <-done:
				fmt.Println("job", j, "canceled")
				return
			case <-time.After(time.Duration(j*10) * time.Millisecond):
				result <- j * j
			}
		}(job)
	}

	// 给工作 goroutine 一点时间，然后发出关闭信号
	time.Sleep(60 * time.Millisecond)
	close(done)

	// 收集已完成的任务结果
	close(result)
	var sum int
	for r := range result {
		sum += r
	}
	fmt.Println("graceful shutdown result sum:", sum)
}

// ==================== nil channel 禁用 select 分支 ====================
//
// 将某个 case 的 channel 设为 nil，可以让 select 忽略该分支。
// 常用于根据状态动态启用/禁用某些 channel。

func NilChannelSelectDemo() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "only ch1 active"
	}()

	// 第一阶段：只监听 ch1
	ch2 = nil
	select {
	case msg := <-ch1:
		fmt.Println("phase1:", msg)
	case msg := <-ch2:
		fmt.Println("phase1:", msg) // 永远不会执行
	}
}
