package safety

import (
	"sync"
	"sync/atomic"
)

type SafeCounter struct {
	mu     sync.RWMutex
	counts map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{
		counts: make(map[string]int),
	}
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	c.counts[key]++
	c.mu.Unlock()
}

func (c *SafeCounter) Get(key string) int {
	c.mu.RLock()
	v := c.counts[key]
	c.mu.RUnlock()
	return v
}

func (c *SafeCounter) Snapshot() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

func ParallelIncrements(c *SafeCounter, key string, goroutines int, each int) {
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				c.Inc(key)
			}
		}()
	}
	wg.Wait()
}

type AtomicCounter struct {
	v atomic.Int64
}

func (c *AtomicCounter) Inc() {
	c.v.Add(1)
}

func (c *AtomicCounter) Value() int64 {
	return c.v.Load()
}
