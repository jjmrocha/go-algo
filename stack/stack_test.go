package stack

import "testing"

func TestNew(t *testing.T) {
	// when
	s := New[int]()
	// then
	if s == nil {
		t.Fatalf("New returned nil")
	}

	if s.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", s.Len())
	}

	if !s.Empty() {
		t.Fatalf("Expected empty stack")
	}
}

func TestPush(t *testing.T) {
	// given
	s := New[int]()
	// when
	s.Push(1)
	s.Push(2)
	s.Push(3)
	// then
	if s.Len() != 3 {
		t.Fatalf("Expected size 3, got %d", s.Len())
	}

	if s.Empty() {
		t.Fatalf("Expected non-empty stack")
	}
}

func TestPop(t *testing.T) {
	// given
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)
	// when
	got, ok := s.Pop()
	// then
	if !ok {
		t.Fatalf("Pop returned false, expected true")
	}

	if got != 3 {
		t.Fatalf("Expected 3, got %d", got)
	}

	if s.Len() != 2 {
		t.Fatalf("Expected size 2 after pop, got %d", s.Len())
	}
}

func TestPopEmpty(t *testing.T) {
	// given
	s := New[int]()
	// when
	got, ok := s.Pop()
	// then
	if ok {
		t.Fatalf("Pop on empty stack should return false")
	}

	if got != 0 {
		t.Fatalf("Pop on empty stack should return zero value, got %d", got)
	}
}

func TestPeek(t *testing.T) {
	// given
	s := New[string]()
	s.Push("a")
	s.Push("b")
	// when
	got, ok := s.Peek()
	// then
	if !ok {
		t.Fatalf("Peek returned false on non-empty stack")
	}

	if got != "b" {
		t.Fatalf("Expected \"b\", got %q", got)
	}

	if s.Len() != 2 {
		t.Fatalf("Peek must not change size, got %d", s.Len())
	}
}

func TestPeekEmpty(t *testing.T) {
	// given
	s := New[string]()
	// when
	got, ok := s.Peek()
	// then
	if ok {
		t.Fatalf("Peek on empty stack should return false")
	}

	if got != "" {
		t.Fatalf("Peek on empty stack should return zero value, got %q", got)
	}
}

func TestGenericType(t *testing.T) {
	// given
	s := New[string]()
	s.Push("hello")
	s.Push("world")
	// when
	got, ok := s.Pop()
	// then
	if !ok {
		t.Fatalf("Pop returned false on non-empty stack")
	}

	if got != "world" {
		t.Fatalf("Expected \"world\", got %q", got)
	}
}
