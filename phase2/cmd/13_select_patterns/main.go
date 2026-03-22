package main

import (
	"fmt"
	"time"

	"go-tutorial/phase2/internal/concurrency"
)

func main() {
	fmt.Println("===== select 随机选择 =====")
	concurrency.SelectRandomDemo()

	fmt.Println("\n===== select with default（非阻塞）=====")
	ch := make(chan string)
	concurrency.NonBlockingReceive(ch)
	concurrency.NonBlockingSend(ch, "hello")

	fmt.Println("\n===== ticker demo =====")
	concurrency.TickerDemo(120 * time.Millisecond)

	fmt.Println("\n===== graceful shutdown =====")
	concurrency.GracefulShutdown([]int{1, 2, 3, 4, 5})

	fmt.Println("\n===== nil channel 禁用分支 =====")
	concurrency.NilChannelSelectDemo()

	fmt.Println("\n===== 超时接收（已有示例）=====")
	slow := make(chan int)
	go func() {
		time.Sleep(50 * time.Millisecond)
		slow <- 99
		close(slow)
	}()

	_, err := concurrency.ReceiveWithTimeout(slow, 10*time.Millisecond)
	if err != nil {
		fmt.Println("timeout as expected:", err)
	}
}
