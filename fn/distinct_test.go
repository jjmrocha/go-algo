package fn

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistinct(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "nil input", input: nil, expected: []int{}},
		{name: "empty input", input: []int{}, expected: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "all duplicates", input: []int{7, 7, 7}, expected: []int{7}},
		{name: "mixed", input: []int{1, 2, 1, 3, 2, 4}, expected: []int{1, 2, 3, 4}},
		{name: "single element", input: []int{5}, expected: []int{5}},
		{name: "duplicates at end", input: []int{1, 2, 3, 3}, expected: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := Distinct(input)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("preserves first occurrence", func(t *testing.T) {
		// given
		input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
		// when
		result := Distinct(input)
		// then
		expected := []int{3, 1, 4, 5, 9, 2, 6}
		assert.Equal(t, expected, result)
	})
}

func TestDistinctSeq(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty", input: []int{}, expected: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "all duplicates", input: []int{7, 7, 7}, expected: []int{7}},
		{name: "mixed", input: []int{1, 2, 1, 3, 2, 4}, expected: []int{1, 2, 3, 4}},
		{name: "single element", input: []int{5}, expected: []int{5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := slices.Collect(DistinctSeq(slices.Values(input)))
			// then
			assert.ElementsMatch(t, tt.expected, result)
		})
	}

	t.Run("multiple iterations", func(t *testing.T) {
		// given
		input := []int{1, 2, 1, 3}
		seq := DistinctSeq(slices.Values(input))
		// when
		first := slices.Collect(seq)
		second := slices.Collect(seq)
		// then
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, first)
		assert.Equal(t, expected, second, "state leaked across iterations")
	})
}
