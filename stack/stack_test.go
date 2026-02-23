package stack

import "testing"

func TestNew(t *testing.T) {
	s := New[int]()
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.Size() != 0 {
		t.Errorf("expected size 0, got %d", s.Size())
	}
	if !s.Empty() {
		t.Error("expected empty stack")
	}
}

func TestPush(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Size() != 3 {
		t.Errorf("expected size 3, got %d", s.Size())
	}
	if s.Empty() {
		t.Error("expected non-empty stack")
	}
}

func TestPop(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	for _, want := range []int{3, 2, 1} {
		got, ok := s.Pop()
		if !ok {
			t.Fatalf("Pop returned false, expected true")
		}
		if got != want {
			t.Errorf("Pop() = %d, want %d", got, want)
		}
	}

	if s.Size() != 0 {
		t.Errorf("expected size 0 after all pops, got %d", s.Size())
	}
}

func TestPopEmpty(t *testing.T) {
	s := New[int]()

	got, ok := s.Pop()
	if ok {
		t.Error("Pop on empty stack should return false")
	}
	if got != 0 {
		t.Errorf("Pop on empty stack should return zero value, got %d", got)
	}
}

func TestPeek(t *testing.T) {
	s := New[string]()
	s.Push("a")
	s.Push("b")

	got, ok := s.Peek()
	if !ok {
		t.Fatal("Peek returned false on non-empty stack")
	}
	if got != "b" {
		t.Errorf("Peek() = %q, want %q", got, "b")
	}
	// Size must be unchanged after Peek
	if s.Size() != 2 {
		t.Errorf("Peek must not change size, got %d", s.Size())
	}
}

func TestPeekEmpty(t *testing.T) {
	s := New[string]()

	got, ok := s.Peek()
	if ok {
		t.Error("Peek on empty stack should return false")
	}
	if got != "" {
		t.Errorf("Peek on empty stack should return zero value, got %q", got)
	}
}

func TestLIFOOrder(t *testing.T) {
	s := New[int]()
	input := []int{10, 20, 30, 40, 50}
	for _, v := range input {
		s.Push(v)
	}

	for i := len(input) - 1; i >= 0; i-- {
		got, ok := s.Pop()
		if !ok {
			t.Fatalf("Pop returned false at index %d", i)
		}
		if got != input[i] {
			t.Errorf("LIFO order violation: got %d, want %d", got, input[i])
		}
	}
}

func TestGenericString(t *testing.T) {
	s := New[string]()
	s.Push("hello")
	s.Push("world")

	got, ok := s.Pop()
	if !ok || got != "world" {
		t.Errorf("Pop() = %q, %v; want %q, true", got, ok, "world")
	}
}
