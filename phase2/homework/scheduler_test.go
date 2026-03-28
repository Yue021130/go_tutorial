package main

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestScheduler_Run(t *testing.T) {
	scheduler := NewScheduler(2)

	var counter int32
	var tasks []Task
	for i := 0; i < 5; i++ {
		tasks = append(tasks, NewSimpleTask(fmt.Sprintf("task-%d", i), func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}))
	}

	results := scheduler.Run(context.Background(), tasks)
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}
	if counter != 5 {
		t.Fatalf("expected counter=5, got %d", counter)
	}

	for _, r := range results {
		if r.Err != nil {
			t.Fatalf("unexpected error for %s: %v", r.TaskID, r.Err)
		}
	}
}

func TestScheduler_RunWithCancellation(t *testing.T) {
	scheduler := NewScheduler(2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var tasks []Task
	for i := 0; i < 10; i++ {
		tasks = append(tasks, NewSimpleTask(fmt.Sprintf("task-%d", i), func(ctx context.Context) error {
			select {
			case <-time.After(100 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}))
	}

	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	results := scheduler.Run(ctx, tasks)
	if len(results) != 10 {
		t.Fatalf("expected 10 results, got %d", len(results))
	}

	var canceled int
	for _, r := range results {
		if r.Err != nil {
			canceled++
		}
	}
	if canceled == 0 {
		t.Fatal("expected at least one canceled task")
	}
}

func TestScheduler_RunWithAggregation(t *testing.T) {
	scheduler := NewScheduler(2)

	tasks := []Task{
		NewSimpleTask("good", func(ctx context.Context) error { return nil }),
		NewSimpleTask("bad", func(ctx context.Context) error { return errors.New("failed") }),
	}

	results, err := scheduler.RunWithAggregation(context.Background(), tasks)
	if err == nil {
		t.Fatal("expected aggregated error")
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestScheduler_EmptyTasks(t *testing.T) {
	scheduler := NewScheduler(2)
	results := scheduler.Run(context.Background(), nil)
	if results != nil {
		t.Fatal("expected nil results for empty tasks")
	}
}
