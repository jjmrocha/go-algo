package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// when
	result := New[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Len())
	assert.True(t, result.Empty())
}

func TestEnqueue(t *testing.T) {
	// given
	q := New[int]()
	// when
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	// then
	assert.Equal(t, int64(3), q.Len())
	assert.False(t, q.Empty())
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
		assert.True(t, ok)
		assert.Equal(t, 10, result)
		assert.Equal(t, int64(1), q.Len())
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := New[int]()
		// when
		result, ok := q.Dequeue()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
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
		assert.True(t, ok)
		assert.Equal(t, expected, result)
	}
}

func TestQueueEmptyAfterDrain(t *testing.T) {
	// given
	q := New[string]()
	q.Enqueue("a")
	q.Dequeue()
	// when / then
	assert.True(t, q.Empty())
	assert.Equal(t, int64(0), q.Len())
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
	assert.True(t, ok)
	assert.Equal(t, 2, result)
}
