package queue

import "testing"

func TestNewQueueIsEmpty(t *testing.T) {
	// when
	q := New[int]()
	// then
	if q == nil {
		t.Fatalf("New returned nil")
	}

	if q.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", q.Len())
	}

	if !q.Empty() {
		t.Fatalf("Expected empty queue")
	}
}

func TestEnqueue(t *testing.T) {
	// given
	q := New[int]()
	// when
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// then
	if q.Len() != 3 {
		t.Fatalf("Expected size 3, got %d", q.Len())
	}

	if q.Empty() {
		t.Fatalf("Expected non-empty queue")
	}
}

func TestDequeue(t *testing.T) {
	// given
	q := New[int]()
	q.Enqueue(10)
	q.Enqueue(20)
	// when
	got, ok := q.Dequeue()
	// then
	if !ok {
		t.Fatalf("Dequeue returned false on non-empty queue")
	}

	if got != 10 {
		t.Fatalf("Expected 10, got %d", got)
	}

	if q.Len() != 1 {
		t.Fatalf("Expected size 1 after dequeue, got %d", q.Len())
	}
}

func TestDequeueEmpty(t *testing.T) {
	// given
	q := New[int]()
	// when
	got, ok := q.Dequeue()
	// then
	if ok {
		t.Fatalf("Dequeue on empty queue should return false")
	}

	if got != 0 {
		t.Fatalf("Dequeue on empty queue should return zero value, got %d", got)
	}
}

func TestQueueFIFOOrdering(t *testing.T) {
	// given
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// when / then
	for _, want := range []int{1, 2, 3} {
		got, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue returned false, expected true")
		}
		if got != want {
			t.Fatalf("Expected %d, got %d", want, got)
		}
	}
}

func TestQueueEmptyAfterDrain(t *testing.T) {
	// given
	q := New[string]()
	q.Enqueue("a")
	q.Dequeue()
	// when / then
	if !q.Empty() {
		t.Fatalf("Queue should be empty after draining all elements")
	}

	if q.Len() != 0 {
		t.Fatalf("Expected size 0 after drain, got %d", q.Len())
	}
}

func TestQueueRefillAfterDrain(t *testing.T) {
	// given
	q := New[int]()
	q.Enqueue(1)
	q.Dequeue()
	// when
	q.Enqueue(2)
	got, ok := q.Dequeue()
	// then
	if !ok {
		t.Fatalf("Dequeue returned false after refill")
	}

	if got != 2 {
		t.Fatalf("Expected 2 after refill, got %d", got)
	}
}
