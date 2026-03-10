package main

import (
	"fmt"

	"go-tutorial/phase2/internal/safety"
)

func main() {
	c := safety.NewSafeCounter()
	safety.ParallelIncrements(c, "orders", 8, 1000)
	fmt.Println("safe counter orders:", c.Get("orders"))

	var ac safety.AtomicCounter
	for i := 0; i < 5; i++ {
		ac.Inc()
	}
	fmt.Println("atomic counter:", ac.Value())
}
