package sorting

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		// given
		arr := []int{}
		// when
		result := Shuffle(arr)
		// then
		assert.NoError(t, result)
		assert.Equal(t, []int{}, arr)
	})

	t.Run("single element", func(t *testing.T) {
		// given
		arr := []int{42}
		// when
		result := Shuffle(arr)
		// then
		assert.NoError(t, result)
		assert.Equal(t, []int{42}, arr)
	})

	t.Run("preserves all elements", func(t *testing.T) {
		// given
		arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		expected := slices.Clone(arr)
		// when
		result := Shuffle(arr)
		// then
		assert.NoError(t, result)
		assert.ElementsMatch(t, expected, arr)
	})
}
