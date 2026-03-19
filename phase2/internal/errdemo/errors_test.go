package errdemo

import "testing"

func TestLoadUserDisplayNameSuccess(t *testing.T) {
	got, err := LoadUserDisplayName(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ALICE" {
		t.Fatalf("expected ALICE, got %s", got)
	}
}

func TestInvalidIDClassification(t *testing.T) {
	_, err := LoadUserDisplayName(0)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !IsInvalidID(err) {
		t.Fatalf("expected invalid id classification")
	}
}

func TestNotFoundClassification(t *testing.T) {
	_, err := LoadUserDisplayName(999)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected not found classification")
	}
}
