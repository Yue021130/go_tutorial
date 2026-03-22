package main

import (
	"fmt"

	"go-tutorial/phase2/internal/safety"
)

func main() {
	fmt.Println("===== sync.RWMutex + map =====")
	c := safety.NewSafeCounter()
	safety.ParallelIncrements(c, "orders", 8, 1000)
	fmt.Println("safe counter orders:", c.Get("orders"))

	fmt.Println("\n===== sync.Mutex =====")
	mc := &safety.MutexCounter{}
	safety.MutexParallelIncrements(mc, 8, 1000)
	fmt.Println("mutex counter value:", mc.Value())

	fmt.Println("\n===== atomic counter =====")
	var ac safety.AtomicCounter
	for i := 0; i < 5; i++ {
		ac.Inc()
	}
	fmt.Println("atomic counter:", ac.Value())

	fmt.Println("\n===== sync.Once =====")
	s1 := safety.GetSingleton()
	s2 := safety.GetSingleton()
	fmt.Printf("singleton1 == singleton2 ? %v, value=%s\n", s1 == s2, s1.GetValue())

	fmt.Println("\n===== sync.Pool =====")
	safety.UseBufferPool()

	fmt.Println("\n===== sync.Map =====")
	safety.SyncMapDemo()

	fmt.Println("\n===== race detector 提示 =====")
	fmt.Println("运行 `go test -race ./...` 或 `go run -race ./cmd/15_sync_and_concurrency_safety` 检测数据竞争")
}
