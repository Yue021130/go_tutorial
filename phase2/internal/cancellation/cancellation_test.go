package cancellation

import (
	"context"
	"testing"
	"time"
)

func TestRunJobsStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	jobs := []int{1, 2, 3, 4, 5, 6}
	got := RunJobs(ctx, jobs, 10*time.Millisecond)

	if len(got) >= len(jobs) {
		t.Fatalf("expected early stop, got all jobs result: %v", got)
	}
	if len(got) == 0 {
		t.Fatalf("expected at least one completed job")
	}
}

func TestTickUntilCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := TickUntilCancel(ctx, 5*time.Millisecond)
	time.Sleep(16 * time.Millisecond)
	cancel()

	count := 0
	for range ch {
		count++
	}
	if count == 0 {
		t.Fatalf("expected some ticks before cancellation")
	}
}
