package stack

import "testing"

func TestStack(t *testing.T) {
	s := New[int]()
	
	if !s.IsEmpty() {
		t.Error("New stack should be empty")
	}
	
	s.Push(1)
	s.Push(2)
	s.Push(3)
	
	if s.Size() != 3 {
		t.Errorf("Expected size 3, got %d", s.Size())
	}
	
	val, ok := s.Peek()
	if !ok || val != 3 {
		t.Errorf("Expected peek to return 3, got %d", val)
	}
	
	val, ok = s.Pop()
	if !ok || val != 3 {
		t.Errorf("Expected pop to return 3, got %d", val)
	}
	
	if s.Size() != 2 {
		t.Errorf("Expected size 2, got %d", s.Size())
	}
}
