package main

import (
	"fmt"

	"go-tutorial/phase2/internal/concurrency"
)

func drain(ch <-chan int) []int {
	out := make([]int, 0)
	for n := range ch {
		out = append(out, n)
	}
	return out
}

func main() {
	fmt.Println("===== unbuffered channel =====")
	concurrency.DemonstrateUnbuffered()

	fmt.Println("\n===== buffered channel =====")
	concurrency.DemonstrateBuffered()

	fmt.Println("\n===== channel close + range =====")
	for n := range concurrency.GenerateWithClose(3) {
		fmt.Println("received:", n)
	}

	fmt.Println("\n===== sync.WaitGroup =====")
	concurrency.WaitGroupExample()

	fmt.Println("\n===== pipeline: generate -> square =====")
	pipelineOut := drain(concurrency.Square(concurrency.Generate(1, 2, 3, 4)))
	fmt.Println("pipeline squares:", pipelineOut)

	fmt.Println("\n===== fan-in: 合并两个流 =====")
	a := concurrency.Square(concurrency.Generate(1, 3, 5))
	b := concurrency.Square(concurrency.Generate(2, 4, 6))
	merged := drain(concurrency.FanIn(a, b))
	fmt.Println("fan-in merged squares:", merged)

	fmt.Println("\n===== worker pool =====")
	poolOut := concurrency.WorkerPoolSquare([]int{1, 2, 3, 4, 5, 6}, 3)
	fmt.Println("worker pool squares(sorted):", poolOut)
}
