package concurrency

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrReceiveTimeout = errors.New("receive timeout")

func Generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func FanIn(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for _, ch := range chans {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for n := range c {
				out <- n
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func ReceiveWithTimeout(ch <-chan int, timeout time.Duration) (int, error) {
	select {
	case n, ok := <-ch:
		if !ok {
			return 0, errors.New("channel closed")
		}
		return n, nil
	case <-time.After(timeout):
		return 0, ErrReceiveTimeout
	}
}

func WorkerPoolSquare(nums []int, workers int) []int {
	if workers <= 0 {
		workers = 1
	}

	jobs := make(chan int)
	results := make(chan int, len(nums))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range jobs {
				results <- n * n
			}
		}()
	}

	go func() {
		for _, n := range nums {
			jobs <- n
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	out := make([]int, 0, len(nums))
	for r := range results {
		out = append(out, r)
	}
	sort.Ints(out)
	return out
}
