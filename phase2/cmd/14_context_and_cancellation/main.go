package main

import (
	"context"
	"fmt"
	"time"

	"go-tutorial/phase2/internal/cancellation"
)

func main() {
	fmt.Println("===== context 超时取消 =====")
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Millisecond)
	defer cancel()

	done := cancellation.RunJobs(ctx, jobs, 10*time.Millisecond)
	fmt.Println("finished jobs result:", done)
	fmt.Println("planned jobs:", jobs)

	fmt.Println("\n===== context.WithCancel =====")
	tickCtx, stop := context.WithCancel(context.Background())
	tickCh := cancellation.TickUntilCancel(tickCtx, 8*time.Millisecond)
	time.Sleep(25 * time.Millisecond)
	stop()

	ticks := make([]int, 0)
	for n := range tickCh {
		ticks = append(ticks, n)
	}
	fmt.Println("ticks before cancel:", ticks)

	fmt.Println("\n===== context.WithValue =====")
	ctx = context.WithValue(context.Background(), cancellation.RequestIDKey, "req-12345")
	ctx = context.WithValue(ctx, cancellation.UserIDKey, 42)
	fmt.Println(cancellation.ProcessRequest(ctx))

	fmt.Println("\n===== context.Deadline =====")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	cancellation.CheckDeadline(timeoutCtx)

	fmt.Println("\n===== 级联取消 =====")
	parent, cancelParent := context.WithCancel(context.Background())
	child1, _ := context.WithCancel(parent)
	child2, _ := context.WithTimeout(parent, 200*time.Millisecond)

	go func() {
		<-child1.Done()
		fmt.Println("child1 canceled, reason:", context.Cause(child1))
	}()
	go func() {
		<-child2.Done()
		fmt.Println("child2 canceled, reason:", context.Cause(child2))
	}()

	time.Sleep(10 * time.Millisecond)
	cancelParent()
	time.Sleep(20 * time.Millisecond)
	fmt.Println("parent canceled, all children notified")
}
