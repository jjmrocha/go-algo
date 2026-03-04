package queue

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
	t.Run("non-empty queue", func(t *testing.T) {
		// given
		q := New[int]()
		q.Enqueue(10)
		q.Enqueue(20)
		// when
		result, ok := q.Dequeue()
		// then
		if !ok {
			t.Fatalf("Dequeue returned false on non-empty queue")
		}

		if result != 10 {
			t.Fatalf("Expected 10, got %d", result)
		}

		if q.Len() != 1 {
			t.Fatalf("Expected size 1 after dequeue, got %d", q.Len())
		}
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := New[int]()
		// when
		result, ok := q.Dequeue()
		// then
		if ok {
			t.Fatalf("Dequeue on empty queue should return false")
		}

		if result != 0 {
			t.Fatalf("Dequeue on empty queue should return zero value, got %d", result)
		}
	})
}

func TestQueueFIFOOrdering(t *testing.T) {
	// given
	q := New[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// when / then
	for _, expected := range []int{1, 2, 3} {
		result, ok := q.Dequeue()
		if !ok {
			t.Fatalf("Dequeue returned false, expected true")
		}
		if result != expected {
			t.Fatalf("Expected %d, got %d", expected, result)
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
	result, ok := q.Dequeue()
	// then
	if !ok {
		t.Fatalf("Dequeue returned false after refill")
	}

	if result != 2 {
		t.Fatalf("Expected 2 after refill, got %d", result)
	}
}
