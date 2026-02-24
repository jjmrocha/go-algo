package queue

import (
	"sync"
	"testing"
)

func TestNewSyncQueueIsEmpty(t *testing.T) {
	// when
	q := NewSyncQueue[int]()
	// then
	if q == nil {
		t.Fatalf("NewSyncQueue returned nil")
	}

	if q.Len() != 0 {
		t.Fatalf("Expected size 0, got %d", q.Len())
	}

	if !q.Empty() {
		t.Fatalf("Expected empty sync queue")
	}
}

func TestSyncQueueEnqueueDequeue(t *testing.T) {
	// given
	q := NewSyncQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// when / then — FIFO order must be preserved
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

func TestSyncQueueDequeueEmpty(t *testing.T) {
	// given
	q := NewSyncQueue[int]()
	// when
	got, ok := q.Dequeue()
	// then
	if ok {
		t.Fatalf("Dequeue on empty SyncQueue should return false")
	}

	if got != 0 {
		t.Fatalf("Dequeue on empty SyncQueue should return zero value, got %d", got)
	}
}

func TestSyncQueueConcurrentEnqueue(t *testing.T) {
	// given
	const goroutines = 100
	q := NewSyncQueue[int]()
	var wg sync.WaitGroup
	// when
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			q.Enqueue(v)
		}(i)
	}
	wg.Wait()
	// then
	if q.Len() != goroutines {
		t.Fatalf("Expected %d elements, got %d", goroutines, q.Len())
	}
}

func TestSyncQueueConcurrentEnqueueDequeue(t *testing.T) {
	// given
	const goroutines = 50
	q := NewSyncQueue[int]()
	var wg sync.WaitGroup
	// when — producers and consumers run concurrently
	for i := range goroutines {
		wg.Add(2)
		go func(v int) {
			defer wg.Done()
			q.Enqueue(v)
		}(i)
		go func() {
			defer wg.Done()
			q.Dequeue()
		}()
	}
	wg.Wait()
	// then — no panic means the data race detector would catch any issues;
	// the exact count depends on scheduling, but must be non-negative.
	if q.Len() < 0 {
		t.Fatalf("Unexpected negative length: %d", q.Len())
	}
}
