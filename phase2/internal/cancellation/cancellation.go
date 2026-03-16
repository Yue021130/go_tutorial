package cancellation

import (
	"context"
	"time"
)

// RunJobs 在 ctx 取消时停止处理后续任务。
func RunJobs(ctx context.Context, jobs []int, perJob time.Duration) []int {
	out := make([]int, 0, len(jobs))
	for _, job := range jobs {
		select {
		case <-ctx.Done():
			return out
		default:
		}

		time.Sleep(perJob)
		out = append(out, job*job)
	}
	return out
}

// TickUntilCancel 周期发出计数，直到 ctx 取消。
func TickUntilCancel(ctx context.Context, interval time.Duration) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		count := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count++
				out <- count
			}
		}
	}()
	return out
}
