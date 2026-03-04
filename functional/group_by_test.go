package functional

import (
	"maps"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupBy(t *testing.T) {
	parity := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}

	t.Run("nil input", func(t *testing.T) {
		// when
		result := GroupBy([]int(nil), parity)
		// then
		assert.Empty(t, result)
	})

	t.Run("empty input", func(t *testing.T) {
		// when
		result := GroupBy([]int{}, parity)
		// then
		assert.Empty(t, result)
	})

	t.Run("all same key", func(t *testing.T) {
		// given
		input := []int{2, 4, 6}
		// when
		result := GroupBy(input, parity)
		// then
		assert.Len(t, result, 1)
		assert.Equal(t, []int{2, 4, 6}, result["even"])
	})

	t.Run("all different keys", func(t *testing.T) {
		// given
		input := []int{1, 2, 3}
		// when
		result := GroupBy(input, func(v int) int { return v })
		// then
		assert.Len(t, result, 3)
		for _, v := range []int{1, 2, 3} {
			assert.Equal(t, []int{v}, result[v])
		}
	})

	t.Run("mixed — order within group preserved", func(t *testing.T) {
		// given
		input := []int{1, 3, 2, 5, 4}
		// when
		result := GroupBy(input, parity)
		// then
		assert.Equal(t, []int{1, 3, 5}, result["odd"])
		assert.Equal(t, []int{2, 4}, result["even"])
	})
}

func TestGroupBySeq(t *testing.T) {
	parity := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}

	t.Run("empty seq", func(t *testing.T) {
		// when
		result := maps.Collect(GroupBySeq(slices.Values([]int{}), parity))
		// then
		assert.Empty(t, result)
	})

	t.Run("all same key", func(t *testing.T) {
		// given
		input := []int{2, 4, 6}
		// when
		result := maps.Collect(GroupBySeq(slices.Values(input), parity))
		// then
		assert.Equal(t, []int{2, 4, 6}, result["even"])
	})

	t.Run("mixed — order within group preserved", func(t *testing.T) {
		// given
		input := []int{1, 2, 3, 4, 5}
		// when
		result := maps.Collect(GroupBySeq(slices.Values(input), parity))
		// then
		assert.Equal(t, []int{1, 3, 5}, result["odd"])
		assert.Equal(t, []int{2, 4}, result["even"])
	})

	t.Run("multiple iterations", func(t *testing.T) {
		// given
		seq := GroupBySeq(slices.Values([]int{1, 2, 3, 4}), parity)
		// when
		first := maps.Collect(seq)
		second := maps.Collect(seq)
		// then
		assert.Equal(t, first["odd"], second["odd"], "odd group: state leaked across iterations")
		assert.Equal(t, first["even"], second["even"], "even group: state leaked across iterations")
	})
}
