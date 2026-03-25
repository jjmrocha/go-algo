package deque

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSyncDeque(t *testing.T) {
	// when
	result := NewSyncDeque[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Len())
	assert.True(t, result.Empty())
}

func TestSyncDequePushFront(t *testing.T) {
	// given
	d := NewSyncDeque[int]()
	// when
	d.PushFront(1)
	d.PushFront(2)
	d.PushFront(3)
	// then
	assert.Equal(t, 3, d.Len())
	assert.False(t, d.Empty())
}

func TestSyncDequePushBack(t *testing.T) {
	// given
	d := NewSyncDeque[int]()
	// when
	d.PushBack(1)
	d.PushBack(2)
	d.PushBack(3)
	// then
	assert.Equal(t, 3, d.Len())
	assert.False(t, d.Empty())
}

func TestSyncDequePopFront(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[int]()
		d.PushBack(10)
		d.PushBack(20)
		// when
		result, ok := d.PopFront()
		// then
		assert.True(t, ok)
		assert.Equal(t, 10, result)
		assert.Equal(t, 1, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[int]()
		// when
		result, ok := d.PopFront()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestSyncDequePopBack(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[int]()
		d.PushBack(10)
		d.PushBack(20)
		// when
		result, ok := d.PopBack()
		// then
		assert.True(t, ok)
		assert.Equal(t, 20, result)
		assert.Equal(t, 1, d.Len())
	})

	t.Run("empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[int]()
		// when
		result, ok := d.PopBack()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestSyncDequePeekFront(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[string]()
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
		d := NewSyncDeque[string]()
		// when
		result, ok := d.PeekFront()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestSyncDequePeekBack(t *testing.T) {
	t.Run("non-empty deque", func(t *testing.T) {
		// given
		d := NewSyncDeque[string]()
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
		d := NewSyncDeque[string]()
		// when
		result, ok := d.PeekBack()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestSyncDequeConcurrentPushFront(t *testing.T) {
	// given
	const goroutines = 100
	d := NewSyncDeque[int]()
	var wg sync.WaitGroup
	// when
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			d.PushFront(v)
		}(i)
	}
	wg.Wait()
	// then
	assert.Equal(t, goroutines, d.Len())
}

func TestSyncDequeConcurrentPushBack(t *testing.T) {
	// given
	const goroutines = 100
	d := NewSyncDeque[int]()
	var wg sync.WaitGroup
	// when
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			d.PushBack(v)
		}(i)
	}
	wg.Wait()
	// then
	assert.Equal(t, goroutines, d.Len())
}

func TestSyncDequeConcurrentMixed(t *testing.T) {
	// given
	const goroutines = 50
	d := NewSyncDeque[int]()
	var wg sync.WaitGroup
	// when — producers push on both ends; consumers pop on both ends concurrently
	for i := range goroutines {
		wg.Add(2)
		go func(v int) {
			defer wg.Done()
			d.PushFront(v)
		}(i)
		go func(v int) {
			defer wg.Done()
			d.PushBack(v)
		}(i)
	}
	for range goroutines {
		wg.Add(2)
		go func() {
			defer wg.Done()
			d.PopFront()
		}()
		go func() {
			defer wg.Done()
			d.PopBack()
		}()
	}
	wg.Wait()
	// then — no panic; the race detector catches any synchronisation failures
	assert.GreaterOrEqual(t, d.Len(), 0)
}
