package deque

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// when
	result := New[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Len())
	assert.True(t, result.Empty())
}

func TestPushFront(t *testing.T) {
	// given
	d := New[int]()
	// when
	d.PushFront(1)
	d.PushFront(2)
	d.PushFront(3)
	// then
	assert.Equal(t, 3, d.Len())
	assert.False(t, d.Empty())
}

func TestPushBack(t *testing.T) {
	// given
	d := New[int]()
	// when
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	// then
	assert.Equal(t, 3, d.Len())
	assert.False(t, d.Empty())
}

func TestPopFront(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(10)
		d.PushBack(20)
		d.PushBack(30)
		// when
		result, ok := d.PopFront()
		// then
		assert.True(t, ok)
		assert.Equal(t, 10, result)
		assert.Equal(t, 2, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := New[int]()
		// when
		result, ok := d.PopFront()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestPopBack(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(10)
		d.PushBack(20)
		d.PushBack(30)
		// when
		result, ok := d.PopBack()
		// then
		assert.True(t, ok)
		assert.Equal(t, 30, result)
		assert.Equal(t, 2, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := New[int]()
		// when
		result, ok := d.PopBack()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestPeekFront(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := New[string]()
		d.PushBack("a")
		d.PushBack("b")
		// when
		result, ok := d.PeekFront()
		// then
		assert.True(t, ok)
		assert.Equal(t, "a", result)
		assert.Equal(t, 2, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := New[string]()
		// when
		result, ok := d.PeekFront()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestPeekBack(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := New[string]()
		d.PushBack("a")
		d.PushBack("b")
		// when
		result, ok := d.PeekBack()
		// then
		assert.True(t, ok)
		assert.Equal(t, "b", result)
		assert.Equal(t, 2, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := New[string]()
		// when
		result, ok := d.PeekBack()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestDequeOrdering(t *testing.T) {
	t.Run("PushBack then PopFront yields FIFO order", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		// when / then
		for _, expected := range []int{1, 2, 3} {
			result, ok := d.PopFront()
			assert.True(t, ok)
			assert.Equal(t, expected, result)
		}
	})

	t.Run("PushFront then PopFront yields LIFO order", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushFront(1)
		d.PushFront(2)
		d.PushFront(3)
		// when / then
		for _, expected := range []int{3, 2, 1} {
			result, ok := d.PopFront()
			assert.True(t, ok)
			assert.Equal(t, expected, result)
		}
	})

	t.Run("PushBack then PopBack yields LIFO order", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		// when / then
		for _, expected := range []int{3, 2, 1} {
			result, ok := d.PopBack()
			assert.True(t, ok)
			assert.Equal(t, expected, result)
		}
	})

	t.Run("mixed PushFront and PushBack preserves front-to-back order", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(2)
		d.PushFront(1)
		d.PushBack(3)
		// when
		var result []int
		for v := range d.Values() {
			result = append(result, v)
		}
		// then — logical order front-to-back: 1, 2, 3
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestDequeGrow(t *testing.T) {
	// given — push beyond defaultCap (8) to trigger two resizes
	d := New[int]()
	const count = 20
	for i := range count {
		d.PushBack(i)
	}
	// when / then — all elements present in FIFO order after resize
	assert.Equal(t, count, d.Len())
	for i := range count {
		result, ok := d.PopFront()
		assert.True(t, ok)
		assert.Equal(t, i, result)
	}
	assert.True(t, d.Empty())
}

func TestDequeGrowFromFront(t *testing.T) {
	// given — PushFront wraps the ring before resize, exercising re-linearisation
	d := New[int]()
	const count = 20
	for i := range count {
		d.PushFront(i)
	}
	// when / then — most-recently pushed item is at the front
	assert.Equal(t, count, d.Len())
	for i := count - 1; i >= 0; i-- {
		result, ok := d.PopFront()
		assert.True(t, ok)
		assert.Equal(t, i, result)
	}
	assert.True(t, d.Empty())
}

func TestDequeShrink(t *testing.T) {
	// given — push 9 items (triggers grow to cap=16), then pop 6 (size=3 < 16/4=4 triggers shrink to cap=8)
	d := New[int]()
	for i := range 9 {
		d.PushBack(i)
	}
	for range 6 {
		d.PopFront()
	}
	// when — deque must preserve correct order and allow further pushes after shrink
	result1, ok1 := d.PopFront()
	result2, ok2 := d.PopFront()
	result3, ok3 := d.PopFront()
	_, ok4 := d.PopFront()
	// then — elements 6, 7, 8 remain after removing 0..5
	assert.True(t, ok1)
	assert.Equal(t, 6, result1)
	assert.True(t, ok2)
	assert.Equal(t, 7, result2)
	assert.True(t, ok3)
	assert.Equal(t, 8, result3)
	assert.False(t, ok4)
}

func TestDequeRefillAfterDrain(t *testing.T) {
	// given — fully drain then refill to verify ring buffer pointers reset correctly
	d := New[int]()
	for i := range 5 {
		d.PushBack(i)
	}
	for range 5 {
		d.PopFront()
	}
	// when
	d.PushBack(10)
	d.PushBack(20)
	result1, ok1 := d.PopFront()
	result2, ok2 := d.PopFront()
	// then
	assert.True(t, ok1)
	assert.Equal(t, 10, result1)
	assert.True(t, ok2)
	assert.Equal(t, 20, result2)
	assert.True(t, d.Empty())
}

func TestDequeValues(t *testing.T) {
	t.Run("iterates front to back without modifying deque", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		// when
		var result []int
		for v := range d.Values() {
			result = append(result, v)
		}
		// then
		assert.Equal(t, []int{1, 2, 3}, result)
		assert.Equal(t, 3, d.Len())
	})

	t.Run("empty deque yields nothing", func(t *testing.T) {
		// given
		d := New[int]()
		// when
		called := false
		for range d.Values() {
			called = true
		}
		// then
		assert.False(t, called)
	})

	t.Run("early termination stops iteration", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		// when
		count := 0
		for range d.Values() {
			count++
			if count == 2 {
				break
			}
		}
		// then — break must not remove elements
		assert.Equal(t, 2, count)
		assert.Equal(t, 3, d.Len())
	})

	t.Run("multiple iterations yield the same sequence", func(t *testing.T) {
		// given
		d := New[int]()
		d.PushBack(1)
		d.PushBack(2)
		d.PushBack(3)
		var first []int
		for v := range d.Values() {
			first = append(first, v)
		}
		// when
		var result []int
		for v := range d.Values() {
			result = append(result, v)
		}
		// then — each range loop gets independent state; deque is unchanged
		assert.Equal(t, first, result)
		assert.Equal(t, 3, d.Len())
	})
}
