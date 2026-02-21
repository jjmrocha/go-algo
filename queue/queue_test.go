package queue

import "testing"

func TestQueue(t *testing.T) {
	q := New[int]()
	
	if !q.IsEmpty() {
		t.Error("New queue should be empty")
	}
	
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	
	if q.Size() != 3 {
		t.Errorf("Expected size 3, got %d", q.Size())
	}
	
	val, ok := q.Peek()
	if !ok || val != 1 {
		t.Errorf("Expected peek to return 1, got %d", val)
	}
	
	val, ok = q.Dequeue()
	if !ok || val != 1 {
		t.Errorf("Expected dequeue to return 1, got %d", val)
	}
	
	if q.Size() != 2 {
		t.Errorf("Expected size 2, got %d", q.Size())
	}
}
