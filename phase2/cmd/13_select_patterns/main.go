package main

import (
	"errors"
	"fmt"
	"time"

	"go-tutorial/phase2/internal/concurrency"
)

func main() {
	slow := make(chan int)
	go func() {
		time.Sleep(50 * time.Millisecond)
		slow <- 99
		close(slow)
	}()

	_, err := concurrency.ReceiveWithTimeout(slow, 10*time.Millisecond)
	if errors.Is(err, concurrency.ErrReceiveTimeout) {
		fmt.Println("select timeout branch triggered")
	}

	n, err := concurrency.ReceiveWithTimeout(slow, 100*time.Millisecond)
	if err != nil {
		fmt.Println("receive error:", err)
		return
	}
	fmt.Println("received with enough timeout:", n)
}
