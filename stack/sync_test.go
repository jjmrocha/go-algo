package stack

import (
	"sync"
	"testing"
)

func TestSyncStackNew(t *testing.T) {
	// when
	result := NewSyncStack[int]()
	// then
	if result == nil {
		t.Fatalf("NewSyncStack returned nil")
	}

	if result.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", result.Len())
	}

	if !result.Empty() {
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
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)
		// when
		result, ok := s.Pop()
		// then
		if !ok {
			t.Fatalf("Pop returned false, expected true")
		}

		if result != 3 {
			t.Fatalf("Expected 3, got %d", result)
		}

		if s.Len() != 2 {
			t.Fatalf("Expected size 2 after pop, got %d", s.Len())
		}
	})

	t.Run("empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[int]()
		// when
		result, ok := s.Pop()
		// then
		if ok {
			t.Fatalf("Pop on empty SyncStack should return false")
		}

		if result != 0 {
			t.Fatalf("Pop on empty SyncStack should return zero value, got %d", result)
		}
	})
}

func TestSyncStackPeek(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[string]()
		s.Push("a")
		s.Push("b")
		// when
		result, ok := s.Peek()
		// then
		if !ok {
			t.Fatalf("Peek returned false on non-empty stack")
		}

		if result != "b" {
			t.Fatalf("Expected \"b\", got %q", result)
		}

		if s.Len() != 2 {
			t.Fatalf("Peek must not change size, got %d", s.Len())
		}
	})

	t.Run("empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[string]()
		// when
		result, ok := s.Peek()
		// then
		if ok {
			t.Fatalf("Peek on empty SyncStack should return false")
		}

		if result != "" {
			t.Fatalf("Peek on empty SyncStack should return zero value, got %q", result)
		}
	})
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
