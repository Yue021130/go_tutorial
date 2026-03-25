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

func TestNilInterface(t *testing.T) {
	var r Repository
	if r != nil {
		t.Fatal("nil interface should be nil")
	}

	var ur *UserRepository
	r = ur
	if r == nil {
		t.Fatal("typed nil assigned to interface should not be nil")
	}

	_, err := r.Find(1)
	if err == nil {
		t.Fatal("expected error for typed nil")
	}
}

func TestStringReader(t *testing.T) {
	r := NewStringReader("hello")
	buf := make([]byte, 3)
	n, err := r.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 3 || string(buf[:n]) != "hel" {
		t.Fatalf("expected 'hel', got %q", string(buf[:n]))
	}
}
