// Package main 提供 Phase 2 实战作业：并发任务调度器。
//
// 功能：
//   - 使用 worker pool 并发执行任务
//   - 支持 context 取消
//   - 支持错误聚合
//   - 支持限流（最大并发数）
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Task 表示一个可执行的任务
type Task interface {
	ID() string
	Execute(ctx context.Context) error
}

// SimpleTask 是一个简单的 Task 实现
type SimpleTask struct {
	id string
	fn func(ctx context.Context) error
}

func NewSimpleTask(id string, fn func(ctx context.Context) error) *SimpleTask {
	return &SimpleTask{id: id, fn: fn}
}

func (t *SimpleTask) ID() string {
	return t.id
}

func (t *SimpleTask) Execute(ctx context.Context) error {
	return t.fn(ctx)
}

// Result 保存任务执行结果
type Result struct {
	TaskID string
	Err    error
}

// Scheduler 并发任务调度器
type Scheduler struct {
	workers int
}

func NewScheduler(workers int) *Scheduler {
	if workers <= 0 {
		workers = 1
	}
	return &Scheduler{workers: workers}
}

// Run 并发执行任务，支持 context 取消
// 返回所有任务的结果，即使某些任务失败也会继续执行其他任务
func (s *Scheduler) Run(ctx context.Context, tasks []Task) []Result {
	if len(tasks) == 0 {
		return nil
	}

	jobs := make(chan Task, len(tasks))
	results := make(chan Result, len(tasks))

	var wg sync.WaitGroup
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				err := task.Execute(ctx)
				results <- Result{TaskID: task.ID(), Err: err}
			}
		}()
	}

	// 发送任务
	go func() {
		for _, task := range tasks {
			jobs <- task
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}
	return out
}

// RunWithAggregation 执行任务并聚合错误
func (s *Scheduler) RunWithAggregation(ctx context.Context, tasks []Task) ([]Result, error) {
	results := s.Run(ctx, tasks)

	var errs []error
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("task %s failed: %w", r.TaskID, r.Err))
		}
	}

	if len(errs) > 0 {
		return results, errors.Join(errs...)
	}
	return results, nil
}
