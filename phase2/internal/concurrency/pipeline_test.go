package concurrency

import (
	"errors"
	"testing"
	"time"
)

func collect(ch <-chan int) []int {
	out := make([]int, 0)
	for n := range ch {
		out = append(out, n)
	}
	return out
}

func TestSquarePipeline(t *testing.T) {
	got := collect(Square(Generate(1, 2, 3)))
	want := []int{1, 4, 9}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got=%v want=%v", i, got, want)
		}
	}
}

func TestReceiveWithTimeout(t *testing.T) {
	ch := make(chan int)
	_, err := ReceiveWithTimeout(ch, 5*time.Millisecond)
	if !errors.Is(err, ErrReceiveTimeout) {
		t.Fatalf("expected ErrReceiveTimeout, got %v", err)
	}
}

func TestWorkerPoolSquare(t *testing.T) {
	got := WorkerPoolSquare([]int{1, 2, 3, 4}, 2)
	want := []int{1, 4, 9, 16}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("index %d mismatch: got=%v want=%v", i, got, want)
		}
	}
}
