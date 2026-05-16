// Program 20_memory_model_gc 演示 Go 内存模型、GC 触发与内存统计观察。
//
// 与 Java 对比：
//   - Java 有 JVM 堆、方法区、栈等明确分区；Go 的内存由运行时统一管理，
//     但逻辑上仍可分为栈（每个 goroutine）和堆（全局）。
//   - Java GC 算法多样（Serial/Parallel/CMS/G1/ZGC/Shenandoah）；
//     Go 当前主流是并发三色标记清除（concurrent tri-color mark-sweep），
//     Go 1.8+ 后 STW（Stop The World）已大幅缩短。
//   - Java 可通过 -Xms/-Xmx 调堆大小；Go 可通过 GOGC 环境变量或 runtime/debug.SetGCPercent
//     调整 GC 触发阈值，默认为 100（堆内存相比上次 GC 后存活对象翻倍时触发）。
//   - Java 有 finalize()（已弃用）；Go 没有析构函数，依赖 GC 回收内存，
//     但 sync.Pool 可用于对象复用。
package main

import (
	"fmt"
	"runtime"

	"go-tutorial/phase4/internal/gc"
)

func main() {
	fmt.Println("=== Go 内存模型与 GC 演示 ===")

	// 打印初始内存状态
	gc.PrintMemStats("初始")

	// 分配一些堆内存（大切片）
	data := make([][]byte, 0, 1000)
	for i := 0; i < 1000; i++ {
		buf := make([]byte, 1024) // 1KB，大概率分配在堆上
		data = append(data, buf)
	}
	gc.PrintMemStats("分配 1000 个 1KB 对象后")

	// 释放引用，让 GC 回收
	data = nil
	gc.ForceGC()
	gc.PrintMemStats("释放并强制 GC 后")

	// 演示 GOGC 调整
	old := gc.SetGCPercent(50)
	fmt.Printf("原 GOGC=%d，已调整为 50（堆内存增长 50%% 即触发 GC）\n", old)
	gc.SetGCPercent(old) // 恢复默认值

	// 打印当前 goroutine 数量，说明栈是 per-goroutine 的
	fmt.Printf("当前 goroutine 数量: %d\n", runtime.NumGoroutine())

	fmt.Println("=== 演示结束 ===")
}
