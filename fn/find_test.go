package fn

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("nil input", func(t *testing.T) {
		// when
		result, ok := Find([]int(nil), isEven)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("empty input", func(t *testing.T) {
		// when
		result, ok := Find([]int{}, isEven)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("no match", func(t *testing.T) {
		// given
		input := []int{1, 3, 5}
		// when
		result, ok := Find(input, isEven)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("first element matches", func(t *testing.T) {
		// given
		input := []int{2, 3, 5}
		// when
		result, ok := Find(input, isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 2, result)
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		// given
		input := []int{1, 4, 6}
		// when
		result, ok := Find(input, isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 4, result)
	})

	t.Run("last element matches", func(t *testing.T) {
		// given
		input := []int{1, 3, 6}
		// when
		result, ok := Find(input, isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 6, result)
	})
}

func TestFindSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("empty seq", func(t *testing.T) {
		// when
		result, ok := FindSeq(slices.Values([]int{}), isEven)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("no match", func(t *testing.T) {
		// given
		input := []int{1, 3, 5}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})

	t.Run("first element matches", func(t *testing.T) {
		// given
		input := []int{2, 3, 5}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 2, result)
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		// given
		input := []int{1, 4, 6}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 4, result)
	})

	t.Run("last element matches", func(t *testing.T) {
		// given
		input := []int{1, 3, 6}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		assert.True(t, ok)
		assert.Equal(t, 6, result)
	})
}
