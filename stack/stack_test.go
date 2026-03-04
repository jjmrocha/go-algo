package stack

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

func TestPush(t *testing.T) {
	// given
	s := New[int]()
	// when
	s.Push(1)
	s.Push(2)
	s.Push(3)
	// then
	assert.Equal(t, int64(3), s.Len())
	assert.False(t, s.Empty())
}

func TestPop(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := New[int]()
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
		s := New[int]()
		// when
		result, ok := s.Pop()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestPeek(t *testing.T) {
	t.Run("non-empty stack", func(t *testing.T) {
		// given
		s := New[string]()
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
		s := New[string]()
		// when
		result, ok := s.Peek()
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})
}

func TestGenericType(t *testing.T) {
	// given
	s := New[string]()
	s.Push("hello")
	s.Push("world")
	// when
	result, ok := s.Pop()
	// then
	assert.True(t, ok)
	assert.Equal(t, "world", result)
}
