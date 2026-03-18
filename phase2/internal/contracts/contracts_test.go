package contracts

import (
	"context"
	"errors"
	"testing"
)

func TestMessageSenderNilNotifier(t *testing.T) {
	s := MessageSender{}
	err := s.Send(context.Background(), "hello")
	if !errors.Is(err, ErrNilNotifier) {
		t.Fatalf("expected ErrNilNotifier, got %v", err)
	}
}

func TestMemoryNotifierNotifyAndLog(t *testing.T) {
	m := &MemoryNotifier{}
	if err := m.Notify(context.Background(), "a"); err != nil {
		t.Fatalf("notify failed: %v", err)
	}
	m.Log(context.Background(), "b")

	if len(m.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(m.Messages))
	}
}

func TestAnyToInt(t *testing.T) {
	n, ok := AnyToInt("123")
	if !ok || n != 123 {
		t.Fatalf("expected 123,true got %d,%t", n, ok)
	}

	_, ok = AnyToInt("abc")
	if ok {
		t.Fatalf("expected false for non-numeric string")
	}
}
