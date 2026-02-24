package stack

import (
	"sync"
	"testing"
)

func TestSyncStackNew(t *testing.T) {
	// when
	s := NewSyncStack[int]()
	// then
	if s == nil {
		t.Fatalf("NewSyncStack returned nil")
	}

	if s.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", s.Len())
	}

	if !s.Empty() {
		t.Fatalf("Expected empty sync stack")
	}
}

func TestSyncStackPush(t *testing.T) {
	// given
	s := NewSyncStack[int]()
	// when
	s.Push(1)
	s.Push(2)
	s.Push(3)
	// then
	if s.Len() != 3 {
		t.Fatalf("Expected size 3, got %d", s.Len())
	}

	if s.Empty() {
		t.Fatalf("Expected non-empty sync stack")
	}
}

func TestSyncStackPop(t *testing.T) {
	// given
	s := NewSyncStack[int]()
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

func TestSyncStackPopEmpty(t *testing.T) {
	// given
	s := NewSyncStack[int]()
	// when
	got, ok := s.Pop()
	// then
	if ok {
		t.Fatalf("Pop on empty SyncStack should return false")
	}

	if got != 0 {
		t.Fatalf("Pop on empty SyncStack should return zero value, got %d", got)
	}
}

func TestSyncStackPeek(t *testing.T) {
	// given
	s := NewSyncStack[string]()
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

func TestSyncStackPeekEmpty(t *testing.T) {
	// given
	s := NewSyncStack[string]()
	// when
	got, ok := s.Peek()
	// then
	if ok {
		t.Fatalf("Peek on empty SyncStack should return false")
	}

	if got != "" {
		t.Fatalf("Peek on empty SyncStack should return zero value, got %q", got)
	}
}

func TestSyncStackConcurrentPush(t *testing.T) {
	// given
	const goroutines = 100
	s := NewSyncStack[int]()
	var wg sync.WaitGroup
	// when
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.Push(v)
		}(i)
	}
	wg.Wait()
	// then
	if s.Len() != goroutines {
		t.Fatalf("Expected %d elements, got %d", goroutines, s.Len())
	}
}
