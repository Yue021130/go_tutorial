package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建 3 个 worker 的调度器
	scheduler := NewScheduler(3)

	// 创建 10 个任务
	var tasks []Task
	for i := 1; i <= 10; i++ {
		id := i
		tasks = append(tasks, NewSimpleTask(fmt.Sprintf("task-%d", id), func(ctx context.Context) error {
			// 模拟工作
			select {
			case <-time.After(20 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}))
	}

	fmt.Println("===== 正常执行 =====")
	results, err := scheduler.RunWithAggregation(context.Background(), tasks)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("all %d tasks completed\n", len(results))
	}

	fmt.Println("\n===== context 取消 =====")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	results, err = scheduler.RunWithAggregation(ctx, tasks)
	if err != nil {
		fmt.Println("expected error:", err)
	}

	fmt.Println("\n===== 部分任务失败 =====")
	tasks = []Task{
		NewSimpleTask("good", func(ctx context.Context) error { return nil }),
		NewSimpleTask("bad", func(ctx context.Context) error { return fmt.Errorf("something wrong") }),
	}
	results, err = scheduler.RunWithAggregation(context.Background(), tasks)
	if err != nil {
		fmt.Println("expected error:", err)
	}
	for _, r := range results {
		fmt.Printf("  %s: %v\n", r.TaskID, r.Err)
	}
}
