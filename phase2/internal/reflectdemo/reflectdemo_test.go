package reflectdemo

import (
	"reflect"
	"testing"
)

func TestModifyByReflect(t *testing.T) {
	u := &User{ID: 1, Name: "Alice"}
	ModifyByReflect(u)

	if u.ID != 999 {
		t.Errorf("expected ID=999, got %d", u.ID)
	}
	if u.Name != "ModifiedByReflect" {
		t.Errorf("expected Name=ModifiedByReflect, got %s", u.Name)
	}
}

func TestCallMethod(t *testing.T) {
	u := User{ID: 1, Name: "Bob"}
	result, err := CallMethod(u)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello, Bob" {
		t.Errorf("expected 'hello, Bob', got %s", result)
	}
}

func TestCompare(t *testing.T) {
	a := []int{1, 2}
	b := []int{1, 2}
	if !reflect.DeepEqual(a, b) {
		t.Error("expected deep equal")
	}
}
