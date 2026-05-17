package tree

import (
	"testing"

	"github.com/jjmrocha/go-algo/sort"
	"github.com/stretchr/testify/assert"
)

func intCmp(a, b int) int {
	if a < b {
		return sort.Before
	}
	if a > b {
		return sort.After
	}
	return sort.Equal
}

func TestNew(t *testing.T) {
	// when
	tr := New[int, int](intCmp)
	// then
	assert.Equal(t, 0, tr.Len())
	assert.True(t, tr.Empty())
}

func TestGet(t *testing.T) {
	t.Run("empty tree returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		result, ok := tr.Get(1)
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})

	t.Run("absent key returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		result, ok := tr.Get(2)
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})

	t.Run("present key returns value and true", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		result, ok := tr.Get(1)
		// then
		assert.True(t, ok)
		assert.Equal(t, "a", result)
	})

	t.Run("returns updated value after overwrite", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(1, "b")
		// when
		result, ok := tr.Get(1)
		// then
		assert.True(t, ok)
		assert.Equal(t, "b", result)
	})
}

func TestContains(t *testing.T) {
	t.Run("false for empty tree", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		result := tr.Contains(1)
		// then
		assert.False(t, result)
	})

	t.Run("false for absent key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.Contains(2)
		// then
		assert.False(t, result)
	})

	t.Run("true for present key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.Contains(1)
		// then
		assert.True(t, result)
	})
}

func TestLen(t *testing.T) {
	t.Run("zero when empty", func(t *testing.T) {
		// given / when
		tr := New[int, int](intCmp)
		// then
		assert.Equal(t, 0, tr.Len())
	})

	t.Run("counts unique keys", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(2, 2)
		tr.Put(3, 3)
		// when
		result := tr.Len()
		// then
		assert.Equal(t, 3, result)
	})

	t.Run("unchanged when overwriting existing key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		tr.Put(1, 99)
		// then
		assert.Equal(t, 1, tr.Len())
	})
}

func TestEmpty(t *testing.T) {
	t.Run("true when new", func(t *testing.T) {
		// given / when
		tr := New[int, int](intCmp)
		// then
		assert.True(t, tr.Empty())
	})

	t.Run("false after put", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		tr.Put(1, 1)
		// then
		assert.False(t, tr.Empty())
	})
}

func TestPut(t *testing.T) {
	t.Run("inserts new key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		tr.Put(42, 99)
		// then
		v, ok := tr.Get(42)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})

	t.Run("overwrites value for existing key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		tr.Put(1, 99)
		// then
		v, ok := tr.Get(1)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})
}

func TestToList(t *testing.T) {
	t.Run("empty tree returns empty slice", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		result := tr.ToList()
		// then
		assert.Empty(t, result)
	})

	t.Run("single element", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1}, result)
	})

	t.Run("ascending insertion returns sorted values", func(t *testing.T) {
		// given — forces left-rotation path; exposes nil-dereference bug in naive traversal
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(2, 2)
		tr.Put(3, 3)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("descending insertion returns sorted values", func(t *testing.T) {
		// given — builds left-heavy tree before LLRBT rotations balance it
		tr := New[int, int](intCmp)
		tr.Put(3, 3)
		tr.Put(2, 2)
		tr.Put(1, 1)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("random insertion order returns sorted values", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, result)
	})

	t.Run("overwritten key appears once with latest value", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(1, 99)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{99}, result)
	})
}
