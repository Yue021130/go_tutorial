// Program 21_gmp_scheduler 演示 Go 的 GMP 调度模型。
//
// GMP 含义：
//   - G（Goroutine）：用户态轻量线程，由 Go 运行时管理。
//   - M（Machine）：操作系统线程，由操作系统调度。
//   - P（Processor）：逻辑处理器，持有本地可运行 goroutine 队列，是 G 和 M 的调度中间层。
//
// GMP 模型文字图：
//
//      +-----+     +-----+     +-----+
//      |  G  |     |  G  |     |  G  |   <- 全局可运行队列（Global Runnable Queue）
//      +-----+     +-----+     +-----+
//         \          |          /
//          v         v         v
//      +-----------------------------+
//      |              P              |   <- 逻辑处理器（默认等于 CPU 核心数）
//      |  +-----------------------+  |
//      |  |   Local Run Queue     |  |   <- 每个 P 的本地队列
//      |  | [G1] [G2] [G3] ...    |  |
//      |  +-----------------------+  |
//      +-----------------------------+
//                   |
//                   v
//      +-----------------------------+
//      |              M              |   <- OS 线程
//      |  (running on CPU core)      |
//      +-----------------------------+
//
// 调度流程：
//   1. 程序启动时，runtime 创建 GOMAXPROCS 个 P。
//   2. 每个 goroutine 创建后，优先放入当前 P 的本地队列。
//   3. M 需要绑定 P 才能执行 G；若本地队列为空，M 会尝试从全局队列或其他 P 偷取 G（work stealing）。
//   4. 当 G 阻塞（如系统调用、channel 等待）时，M 可能与 P 解绑，其他 M 继续执行该 P 的队列。
//
// 与 Java 线程调度对比：
//   - Java 线程是 OS 线程（1:1 模型），创建/切换成本高；goroutine 是用户态轻量协程，初始栈仅 2KB，
//     切换由 runtime 完成，不经过内核，成本极低。
//   - Java 线程调度依赖操作系统调度器；Go 有自有调度器，可在用户态快速切换和窃取任务。
//   - Java 线程池（ThreadPoolExecutor）用于限制并发数；Go 通常直接启动大量 goroutine，
//     通过 channel / sync / context 协调，极少需要手动维护线程池。
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"go-tutorial/phase4/internal/gmp"
)

func main() {
	fmt.Println("=== GMP 调度模型演示 ===")
	gmp.PrintGMPInfo()

	// 演示 1：大量 goroutine 并发执行
	const n = 10
	var wg sync.WaitGroup
	wg.Add(n)

	start := time.Now()
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			// 模拟一点 CPU 工作
			sum := 0
			for j := 0; j < 1e6; j++ {
				sum += j
			}
			fmt.Printf("goroutine %d done, sum=%d\n", id, sum)
		}(i)
	}
	wg.Wait()
	fmt.Printf("%d 个 goroutine 执行耗时: %v\n", n, time.Since(start))

	// 演示 2：调整 GOMAXPROCS
	fmt.Printf("\n当前 GOMAXPROCS=%d\n", runtime.GOMAXPROCS(0))
	old := gmp.GOMAXPROCS(2)
	fmt.Printf("调整为 GOMAXPROCS=2（原值=%d），再运行一次并发任务\n", old)

	wg.Add(n)
	start = time.Now()
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			sum := 0
			for j := 0; j < 1e6; j++ {
				sum += j
			}
			fmt.Printf("(P=2) goroutine %d done, sum=%d\n", id, sum)
		}(i)
	}
	wg.Wait()
	fmt.Printf("(P=2) %d 个 goroutine 执行耗时: %v\n", n, time.Since(start))

	// 恢复默认值
	gmp.GOMAXPROCS(old)
	fmt.Println("=== 演示结束 ===")
}
