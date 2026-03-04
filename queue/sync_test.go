package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncQueue(t *testing.T) {
	// when
	result := NewSyncQueue[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Len())
	assert.True(t, result.Empty())
}

func TestSyncQueueEnqueueDequeue(t *testing.T) {
	// given
	q := NewSyncQueue[int]()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// when / then — FIFO order must be preserved
	for _, expected := range []int{1, 2, 3} {
		result, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, expected, result)
	}
}

func TestSyncQueueDequeueEmpty(t *testing.T) {
	// given
	q := NewSyncQueue[int]()
	// when
	result, ok := q.Dequeue()
	// then
	assert.False(t, ok)
	assert.Equal(t, 0, result)
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
	assert.Equal(t, int64(goroutines), q.Len())
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
	assert.GreaterOrEqual(t, q.Len(), int64(0))
}
