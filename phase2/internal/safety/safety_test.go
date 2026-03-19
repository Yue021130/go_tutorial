package safety

import "testing"

func TestParallelIncrements(t *testing.T) {
	c := NewSafeCounter()
	ParallelIncrements(c, "k", 10, 100)

	if got := c.Get("k"); got != 1000 {
		t.Fatalf("expected 1000, got %d", got)
	}
}

func TestAtomicCounter(t *testing.T) {
	var c AtomicCounter
	for i := 0; i < 7; i++ {
		c.Inc()
	}
	if got := c.Value(); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}
