// Package gc 用于演示 Go 内存模型、GC 触发与对象生命周期观察。
package gc

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

// PrintMemStats 打印当前内存统计信息。
func PrintMemStats(label string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("[%s] HeapAlloc=%.2f MB, HeapObjects=%d, NumGC=%d, LastGC=%s\n",
		label,
		float64(m.HeapAlloc)/(1024*1024),
		m.HeapObjects,
		m.NumGC,
		time.Unix(0, int64(m.LastGC)).Format("15:04:05"),
	)
}

// ForceGC 强制触发一次 GC 并等待其完成。
func ForceGC() {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)
}

// SetGCPercent 调整 GC 触发阈值（GOGC 百分比）。
func SetGCPercent(percent int) int {
	return debug.SetGCPercent(percent)
}
