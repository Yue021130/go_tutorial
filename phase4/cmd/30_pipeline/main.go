// Program 30_pipeline 演示 Pipeline 并发模式。
//
// Pipeline 模式：
//   - 将数据处理拆分为多个阶段（stage），每个阶段通过 channel 连接。
//   - 每个阶段可由一个或多个 goroutine 并行处理。
//   - 适合流式数据处理（如日志清洗、ETL、数据转换）。
//
// 数据流示意：
//
//   [Source] --chan--> [Stage 1] --chan--> [Stage 2] --chan--> [Sink]
//      |                    |                   |                   |
//   生成数据            过滤/转换           聚合/输出           消费结果
//
// 与 Java Stream 对比：
//   - Java Stream 是单线程或并行流（fork/join）的函数式链式调用；
//   - Go Pipeline 基于 channel + goroutine，天然并发，背压由 channel 缓冲控制。
package main

import (
	"fmt"
	"sync"
)

// generator 生成整数序列。
func generator(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// square 对输入通道中的每个数求平方。
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// filter 过滤掉大于 max 的值。
func filter(in <-chan int, max int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n <= max {
				out <- n
			}
		}
	}()
	return out
}

// fanOut 把输入分发给多个 worker，并把结果合并。
func fanOut(in <-chan int, workers int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	worker := func() {
		defer wg.Done()
		for n := range in {
			out <- n * 2 // 模拟耗时处理
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("=== Pipeline 模式演示 ===")

	// 构建处理链：生成 -> 求平方 -> 过滤 -> fan-out 处理 -> 输出
	nums := generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(nums)
	filtered := filter(squared, 50)
	processed := fanOut(filtered, 3)

	for result := range processed {
		fmt.Printf("结果: %d\n", result)
	}

	fmt.Println("=== 演示结束 ===")
}
