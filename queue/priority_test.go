package queue

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intCmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func mustNewPriorityQueueWithCap[T any](t testing.TB, cap int, cmp func(a, b T) int) *PriorityQueue[T] {
	t.Helper()
	q, err := NewPriorityQueueWithCap[T](cap, cmp)
	if err != nil {
		t.Fatalf("NewPriorityQueueWithCap(%d): unexpected error: %v", cap, err)
	}
	return q
}

func TestNewPriorityQueue(t *testing.T) {
	// when
	result := NewPriorityQueue[int](intCmp)
	// then
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Len())
	assert.True(t, result.Empty())
}

func TestNewPriorityQueueWithCap(t *testing.T) {
	t.Run("valid capacity", func(t *testing.T) {
		// when
		result, err := NewPriorityQueueWithCap[int](4, intCmp)
		// then
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.Len())
		assert.True(t, result.Empty())
	})

	t.Run("zero capacity", func(t *testing.T) {
		// when
		result, err := NewPriorityQueueWithCap[int](0, intCmp)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
		assert.Nil(t, result)
	})

	t.Run("negative capacity", func(t *testing.T) {
		// when
		result, err := NewPriorityQueueWithCap[int](-1, intCmp)
		// then
		assert.ErrorIs(t, err, ErrCapacityTooSmall)
		assert.Nil(t, result)
	})
}

func TestPriorityQueueEnqueue(t *testing.T) {
	// given
	q := NewPriorityQueue[int](intCmp)
	// when
	q.Enqueue(5)
	q.Enqueue(3)
	q.Enqueue(7)
	// then
	assert.Equal(t, 3, q.Len())
	assert.False(t, q.Empty())
}

func TestPriorityQueueDequeue(t *testing.T) {
	t.Run("non-empty queue returns minimum", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(5)
		q.Enqueue(3)
		q.Enqueue(7)
		// when
		result, ok := q.Dequeue()
		// then
		assert.True(t, ok)
		assert.Equal(t, 3, result)
		assert.Equal(t, 2, q.Len())
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		// when
		result, ok := q.Dequeue()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("queue is empty after draining", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Dequeue()
		// when
		result, ok := q.Dequeue()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("can re-enqueue after draining", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Dequeue()
		// when
		q.Enqueue(2)
		result, ok := q.Dequeue()
		// then
		assert.True(t, ok)
		assert.Equal(t, 2, result)
	})
}

func TestPriorityQueuePeek(t *testing.T) {
	t.Run("non-empty queue returns minimum", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(5)
		q.Enqueue(3)
		q.Enqueue(7)
		// when
		result, ok := q.Peek()
		// then
		assert.True(t, ok)
		assert.Equal(t, 3, result)
	})

	t.Run("empty queue", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		// when
		result, ok := q.Peek()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("does not remove element", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		// when
		q.Peek()
		// then
		assert.Equal(t, 1, q.Len())
	})

	t.Run("returns same element on repeated calls", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(3)
		q.Enqueue(1)
		q.Enqueue(2)
		// when
		result1, ok1 := q.Peek()
		result2, ok2 := q.Peek()
		// then
		assert.True(t, ok1)
		assert.True(t, ok2)
		assert.Equal(t, result1, result2)
	})
}

func TestPriorityQueueLen(t *testing.T) {
	t.Run("zero when empty", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		// when
		result := q.Len()
		// then
		assert.Equal(t, 0, result)
	})

	t.Run("increases after enqueue", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		// when
		result := q.Len()
		// then
		assert.Equal(t, 1, result)
	})

	t.Run("decreases after dequeue", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Dequeue()
		// when
		result := q.Len()
		// then
		assert.Equal(t, 1, result)
	})
}

func TestPriorityQueueEmpty(t *testing.T) {
	t.Run("initially empty", func(t *testing.T) {
		// given / when
		q := NewPriorityQueue[int](intCmp)
		// then
		assert.True(t, q.Empty())
	})

	t.Run("not empty after enqueue", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		// when
		q.Enqueue(1)
		// then
		assert.False(t, q.Empty())
	})

	t.Run("empty after draining", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		// when
		q.Dequeue()
		// then
		assert.True(t, q.Empty())
	})
}

func TestPriorityQueuePriorityOrdering(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "ascending input",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "descending input",
			input:    []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "random order input",
			input:    []int{3, 1, 4, 1, 5, 9, 2, 6},
			expected: []int{1, 1, 2, 3, 4, 5, 6, 9},
		},
		{
			name:     "duplicate elements",
			input:    []int{2, 2, 2},
			expected: []int{2, 2, 2},
		},
		{
			name:     "single element",
			input:    []int{42},
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			q := NewPriorityQueue[int](intCmp)
			for _, v := range input {
				q.Enqueue(v)
			}
			// when
			var result []int
			for v, ok := q.Dequeue(); ok; v, ok = q.Dequeue() {
				result = append(result, v)
			}
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPriorityQueueDrain(t *testing.T) {
	t.Run("drains queue in priority order", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(3)
		q.Enqueue(1)
		q.Enqueue(2)
		// when
		var result []int
		for v := range q.Drain() {
			result = append(result, v)
		}
		// then
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
		assert.True(t, q.Empty())
	})

	t.Run("empty queue yields nothing", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		// when
		called := false
		for range q.Drain() {
			called = true
		}
		// then
		assert.False(t, called)
	})

	t.Run("early termination", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Enqueue(2)
		q.Enqueue(3)
		// when
		count := 0
		for range q.Drain() {
			count++
			if count == 2 {
				break
			}
		}
		// then
		assert.Equal(t, 2, count)
	})

	t.Run("second iteration is empty after queue is drained", func(t *testing.T) {
		// given
		q := NewPriorityQueue[int](intCmp)
		q.Enqueue(1)
		q.Enqueue(2)
		seq := q.Drain()
		first := slices.Collect(seq)
		// when
		result := slices.Collect(seq)
		// then
		expected := []int{1, 2}
		assert.Equal(t, expected, first)
		assert.Empty(t, result)
	})
}

func TestPriorityQueueGrowsAndShrinks(t *testing.T) {
	// given — start with capacity smaller than the number of elements to force multiple resizes
	const n = 100
	q := mustNewPriorityQueueWithCap[int](t, 4, intCmp)
	for i := n; i > 0; i-- {
		q.Enqueue(i)
	}
	// when / then — heap property must be preserved after all resizes
	for i := range n {
		expected := i + 1
		result, ok := q.Dequeue()
		assert.True(t, ok)
		assert.Equal(t, expected, result)
	}
	assert.True(t, q.Empty())
}
