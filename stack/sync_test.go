package stack

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncStackNew(t *testing.T) {
	// when
	result := NewSyncStack[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result.Len())
	assert.True(t, result.Empty())
}

func TestSyncStackPush(t *testing.T) {
	// given
	s := NewSyncStack[int]()
	// when
	s.Push(1)
	s.Push(2)
	s.Push(3)
	// then
	assert.Equal(t, int64(3), s.Len())
	assert.False(t, s.Empty())
}

func TestSyncStackPop(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[int]()
		s.Push(1)
		s.Push(2)
		s.Push(3)
		// when
		result, ok := s.Pop()
		// then
		assert.True(t, ok)
		assert.Equal(t, 3, result)
		assert.Equal(t, int64(2), s.Len())
	})

	t.Run("empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[int]()
		// when
		result, ok := s.Pop()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestSyncStackPeek(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[string]()
		s.Push("a")
		s.Push("b")
		// when
		result, ok := s.Peek()
		// then
		assert.True(t, ok)
		assert.Equal(t, "b", result)
		assert.Equal(t, int64(2), s.Len())
	})

	t.Run("empty stack", func(t *testing.T) {
		// given
		s := NewSyncStack[string]()
		// when
		result, ok := s.Peek()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestSyncStackConcurrentPush(t *testing.T) {
	// given
	const goroutines = 100
	s := NewSyncStack[int]()
	var wg sync.WaitGroup
	// when
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			s.Push(v)
		}(i)
	}
	wg.Wait()
	// then
	assert.Equal(t, int64(goroutines), s.Len())
}
