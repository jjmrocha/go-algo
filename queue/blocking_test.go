package queue

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// blockWindow is how long a "still blocked" assertion waits before concluding an
// operation is genuinely blocked; unblock assertions use a generous timeout.
const blockWindow = 50 * time.Millisecond

func TestNewBlockingQueue(t *testing.T) {
	t.Run("valid capacity", func(t *testing.T) {
		// when
		result, err := NewBlockingQueue[int](5)
		// then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 5, result.Cap())
		assert.Equal(t, 0, result.Len())
		assert.True(t, result.Empty())
	})

	t.Run("zero capacity returns error", func(t *testing.T) {
		// when
		_, err := NewBlockingQueue[int](0)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
	})

	t.Run("negative capacity returns error", func(t *testing.T) {
		// when
		_, err := NewBlockingQueue[int](-1)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
	})
}

func TestBlockingQueueEnqueueDequeue(t *testing.T) {
	// given
	q, err := NewBlockingQueue[int](3)
	require.NoError(t, err)
	q.Enqueue(10)
	q.Enqueue(20)
	// when
	result := q.Dequeue()
	// then
	assert.Equal(t, 10, result)
	assert.Equal(t, 1, q.Len())
}

func TestBlockingQueueFIFOOrdering(t *testing.T) {
	// given
	q, err := NewBlockingQueue[int](3)
	require.NoError(t, err)
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// when / then — FIFO order must be preserved
	for _, expected := range []int{1, 2, 3} {
		result := q.Dequeue()
		assert.Equal(t, expected, result)
	}
}

func TestBlockingQueueLen(t *testing.T) {
	// given
	q, err := NewBlockingQueue[int](5)
	require.NoError(t, err)
	// when
	q.Enqueue(1)
	q.Enqueue(2)
	// then
	assert.Equal(t, 2, q.Len())
}

func TestBlockingQueueEmpty(t *testing.T) {
	t.Run("empty queue", func(t *testing.T) {
		// given
		q, err := NewBlockingQueue[int](3)
		require.NoError(t, err)
		// when / then
		assert.True(t, q.Empty())
	})

	t.Run("non-empty queue", func(t *testing.T) {
		// given
		q, err := NewBlockingQueue[int](3)
		require.NoError(t, err)
		q.Enqueue(1)
		// when / then
		assert.False(t, q.Empty())
	})
}

func TestBlockingQueueCap(t *testing.T) {
	// given
	q, err := NewBlockingQueue[int](7)
	require.NoError(t, err)
	// when
	result := q.Cap()
	// then
	assert.Equal(t, 7, result)
}

func TestBlockingQueueFull(t *testing.T) {
	t.Run("not full", func(t *testing.T) {
		// given
		q, err := NewBlockingQueue[int](3)
		require.NoError(t, err)
		q.Enqueue(1)
		// when / then
		assert.False(t, q.Full())
	})

	t.Run("full", func(t *testing.T) {
		// given
		q, err := NewBlockingQueue[int](2)
		require.NoError(t, err)
		q.Enqueue(1)
		q.Enqueue(2)
		// when / then
		assert.True(t, q.Full())
	})
}

func TestBlockingQueueConcurrentEnqueue(t *testing.T) {
	// given
	const n = 50
	q, err := NewBlockingQueue[int](n)
	require.NoError(t, err)
	var wg sync.WaitGroup
	// when: concurrent enqueues
	for i := range n {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			q.Enqueue(v)
		}(i)
	}
	wg.Wait()
	// then
	assert.Equal(t, n, q.Len())
}

func TestBlockingQueueEnqueueBlocksWhenFull(t *testing.T) {
	// given: a queue filled to capacity
	q, err := NewBlockingQueue[int](1)
	require.NoError(t, err)
	q.Enqueue(1)

	done := make(chan struct{})
	// when: a further enqueue must block until a slot frees
	go func() {
		q.Enqueue(2)
		close(done)
	}()

	// then: it stays blocked while the queue is full
	select {
	case <-done:
		t.Fatal("Enqueue returned while the queue was full")
	case <-time.After(blockWindow):
	}

	// and: freeing a slot unblocks it
	require.Equal(t, 1, q.Dequeue())
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Enqueue did not unblock after a slot was freed")
	}
	assert.Equal(t, 2, q.Dequeue())
}

func TestBlockingQueueDequeueBlocksWhenEmpty(t *testing.T) {
	// given: an empty queue
	q, err := NewBlockingQueue[int](1)
	require.NoError(t, err)

	result := make(chan int, 1)
	// when: a dequeue must block until a value is available
	go func() {
		result <- q.Dequeue()
	}()

	// then: it stays blocked while the queue is empty
	select {
	case <-result:
		t.Fatal("Dequeue returned while the queue was empty")
	case <-time.After(blockWindow):
	}

	// and: enqueuing a value unblocks it
	q.Enqueue(42)
	select {
	case v := <-result:
		assert.Equal(t, 42, v)
	case <-time.After(time.Second):
		t.Fatal("Dequeue did not unblock after a value was enqueued")
	}
}
