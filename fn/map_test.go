package fn

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		expected []int
	}{
		{
			name:     "nil input",
			input:    nil,
			fn:       func(v int) int { return v * 2 },
			expected: []int{},
		},
		{
			name:     "empty input",
			input:    []int{},
			fn:       func(v int) int { return v * 2 },
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{3},
			fn:       func(v int) int { return v * 2 },
			expected: []int{6},
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3},
			fn:       func(v int) int { return v * 2 },
			expected: []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := Map(input, tt.fn)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("type transform", func(t *testing.T) {
		// given: int → bool: verify T→U type change works
		input := []int{0, 1, 2}
		// when
		result := Map(input, func(v int) bool { return v > 0 })
		// then
		expected := []bool{false, true, true}
		assert.Equal(t, expected, result)
	})

	t.Run("preserves order", func(t *testing.T) {
		// given
		input := []int{5, 3, 1, 4, 2}
		// when
		result := Map(input, func(v int) int { return v * 10 })
		// then
		expected := []int{50, 30, 10, 40, 20}
		assert.Equal(t, expected, result)
	})
}

func TestMapSeq(t *testing.T) {
	double := func(v int) int { return v * 2 }

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty seq", input: []int{}, expected: []int{}},
		{name: "single element", input: []int{3}, expected: []int{6}},
		{name: "multiple elements", input: []int{1, 2, 3}, expected: []int{2, 4, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := slices.Collect(MapSeq(slices.Values(input), double))
			// then
			assert.ElementsMatch(t, tt.expected, result)
		})
	}

	t.Run("type transform", func(t *testing.T) {
		// given: int → bool: verifies the output type differs from the input type
		input := []int{1, 2, 3, 4}
		// when
		result := slices.Collect(MapSeq(slices.Values(input), func(v int) bool { return v%2 == 0 }))
		// then
		expected := []bool{false, true, false, true}
		assert.Equal(t, expected, result)
	})
}
