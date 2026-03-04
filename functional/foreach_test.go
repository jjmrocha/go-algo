package functional

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEach(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		// given
		var collected []int
		// when
		ForEach([]int{1, 2, 3}, func(v int) { collected = append(collected, v) })
		// then
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, collected)
	})

	t.Run("nil input — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEach([]int(nil), func(_ int) { called = true })
		// then
		assert.False(t, called)
	})

	t.Run("empty input — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEach([]int{}, func(_ int) { called = true })
		// then
		assert.False(t, called)
	})
}

func TestForEachSeq(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		// given
		input := []int{1, 2, 3, 4, 5}
		var collected []int
		// when
		ForEachSeq(slices.Values(input), func(v int) { collected = append(collected, v) })
		// then
		assert.Equal(t, input, collected)
	})

	t.Run("empty seq — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEachSeq(slices.Values([]int{}), func(_ int) { called = true })
		// then
		assert.False(t, called)
	})
}
