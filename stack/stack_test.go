package stack

import "testing"

func TestNew(t *testing.T) {
	// when
	result := New[int]()
	// then
	if result == nil {
		t.Fatalf("New returned nil")
	}

	if result.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", result.Len())
	}

	if !result.Empty() {
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
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := New[int]()
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
		s := New[int]()
		// when
		result, ok := s.Pop()
		// then
		if ok {
			t.Fatalf("Pop on empty stack should return false")
		}

		if result != 0 {
			t.Fatalf("Pop on empty stack should return zero value, got %d", result)
		}
	})
}

func TestPeek(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := New[string]()
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
		s := New[string]()
		// when
		result, ok := s.Peek()
		// then
		if ok {
			t.Fatalf("Peek on empty stack should return false")
		}

		if result != "" {
			t.Fatalf("Peek on empty stack should return zero value, got %q", result)
		}
	})
}

func TestGenericType(t *testing.T) {
	// given
	s := New[string]()
	s.Push("hello")
	s.Push("world")
	// when
	result, ok := s.Pop()
	// then
	if !ok {
		t.Fatalf("Pop returned false on non-empty stack")
	}

	if result != "world" {
		t.Fatalf("Expected \"world\", got %q", result)
	}
}
