package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncPriorityQueue(t *testing.T) {
	// when
	result := NewSyncPriorityQueue[int](intCmp)
	// then
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Len())
	assert.True(t, result.Empty())
}

func TestNewSyncPriorityQueueWithCap(t *testing.T) {
	t.Run("valid capacity", func(t *testing.T) {
		// when
		result, err := NewSyncPriorityQueueWithCap[int](4, intCmp)
		// then
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("zero capacity returns error", func(t *testing.T) {
		// when
		result, err := NewSyncPriorityQueueWithCap[int](0, intCmp)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
		assert.Nil(t, result)
	})

	t.Run("negative capacity returns error", func(t *testing.T) {
		// when
		result, err := NewSyncPriorityQueueWithCap[int](-1, intCmp)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
		assert.Nil(t, result)
	})
}

func TestSyncPQueueDequeue(t *testing.T) {
	t.Run("returns elements in priority order", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		q.Enqueue(3)
		q.Enqueue(1)
		q.Enqueue(2)
		// when / then
		for _, expected := range []int{1, 2, 3} {
			result, ok := q.Dequeue()
			assert.True(t, ok)
			assert.Equal(t, expected, result)
		}
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		// when
		result, ok := q.Dequeue()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestSyncPQueuePeek(t *testing.T) {
	t.Run("returns minimum without removing", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		q.Enqueue(5)
		q.Enqueue(1)
		q.Enqueue(3)
		// when
		result, ok := q.Peek()
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, result)
		assert.Equal(t, 3, q.Len())
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		// when
		result, ok := q.Peek()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestSyncPQueueValues(t *testing.T) {
	t.Run("drains queue in priority order", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		q.Enqueue(3)
		q.Enqueue(1)
		q.Enqueue(2)
		// when
		var result []int
		for v := range q.Values() {
			result = append(result, v)
		}
		// then
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
		assert.True(t, q.Empty())
	})

	t.Run("empty queue yields nothing", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		// when
		called := false
		for range q.Values() {
			called = true
		}
		// then
		assert.False(t, called)
	})

	t.Run("early termination", func(t *testing.T) {
		// given
		q := NewSyncPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		// when
		count := 0
		for range q.Values() {
			count++
			if count == 2 {
				break
			}
		}
		// then
		assert.Equal(t, 2, count)
	})
}

func TestSyncPQueueConcurrentEnqueue(t *testing.T) {
	// given
	const goroutines = 100
	q := NewSyncPriorityQueue[int](intCmp)
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
	assert.Equal(t, goroutines, q.Len())
}

func TestSyncPQueueConcurrentEnqueueDequeue(t *testing.T) {
	// given
	const goroutines = 50
	q := NewSyncPriorityQueue[int](intCmp)
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
	// then — no panic; race detector catches data races
	assert.GreaterOrEqual(t, q.Len(), 0)
}

func TestSyncPQueueConcurrentPeek(t *testing.T) {
	// given
	const goroutines = 100
	q := NewSyncPriorityQueue[int](intCmp)
	q.Enqueue(1)
	var wg sync.WaitGroup
	// when — concurrent reads while the queue is non-empty
	for range goroutines {
		wg.Go(func() {
			q.Peek()
		})
	}
	wg.Wait()
	// then — no panic; element is still present
	assert.Equal(t, 1, q.Len())
}
