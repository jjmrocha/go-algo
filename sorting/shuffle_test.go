package sorting

import (
	"errors"
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

func TestShuffleWithRandom(t *testing.T) {
	t.Run("applies the swaps chosen by the source", func(t *testing.T) {
		// given: a source that always selects index 0
		arr := []int{0, 1, 2, 3}
		fn := func(int) (int, error) { return 0, nil }
		// when
		result := ShuffleWithRandom(arr, fn)
		// then: each element i is swapped with index 0 in turn
		assert.NoError(t, result)
		assert.Equal(t, []int{3, 0, 1, 2}, arr)
	})

	t.Run("draws from the range [0,i] at each step", func(t *testing.T) {
		// given
		arr := []int{10, 20, 30, 40}
		var bounds []int
		fn := func(n int) (int, error) {
			bounds = append(bounds, n)
			return 0, nil
		}
		// when
		result := ShuffleWithRandom(arr, fn)
		// then: the source is asked for an index in [0,i] — i.e. n=i+1 — for i in 1..len-1
		assert.NoError(t, result)
		assert.Equal(t, []int{2, 3, 4}, bounds)
	})

	t.Run("propagates an error from the source", func(t *testing.T) {
		// given
		arr := []int{1, 2, 3}
		sentinel := errors.New("rng failure")
		fn := func(int) (int, error) { return 0, sentinel }
		// when
		result := ShuffleWithRandom(arr, fn)
		// then
		assert.ErrorIs(t, result, sentinel)
	})
}
