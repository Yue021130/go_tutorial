package main

import (
	"context"
	"fmt"
	"time"

	"go-tutorial/phase2/internal/cancellation"
)

func main() {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()

	done := cancellation.RunJobs(ctx, jobs, 10*time.Millisecond)
	fmt.Println("finished jobs result:", done)
	fmt.Println("planned jobs:", jobs)

	// Tick 示例：观察取消传播。
	tickCtx, stop := context.WithCancel(context.Background())
	tickCh := cancellation.TickUntilCancel(tickCtx, 8*time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	stop()

	ticks := make([]int, 0)
	for n := range tickCh {
		ticks = append(ticks, n)
	}
	fmt.Println("ticks before cancel:", ticks)
}
