package functional

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFold(t *testing.T) {
	add := func(acc, v int) int { return acc + v }

	tests := []struct {
		name     string
		input    []int
		initial  int
		expected int
	}{
		{name: "nil input returns initial", input: nil, initial: 42, expected: 42},
		{name: "empty input returns initial", input: []int{}, initial: 7, expected: 7},
		{name: "single element", input: []int{5}, initial: 0, expected: 5},
		{name: "sum", input: []int{1, 2, 3, 4}, initial: 0, expected: 10},
		{name: "non-zero initial", input: []int{1, 2, 3}, initial: 10, expected: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			initial := tt.initial
			// when
			result := Fold(input, initial, add)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("left order", func(t *testing.T) {
		// subtraction is non-commutative: verifies left-to-right application
		// given
		input := []int{1, 2, 3}
		// when: ((10 - 1) - 2) - 3 = 4
		result := Fold(input, 10, func(acc, v int) int { return acc - v })
		// then
		assert.Equal(t, 4, result)
	})

	t.Run("string concat", func(t *testing.T) {
		// given
		input := []string{"b", "c", "d"}
		// when
		result := Fold(input, "a", func(acc, v string) string { return acc + v })
		// then
		assert.Equal(t, "abcd", result)
	})
}

func TestFoldSeq(t *testing.T) {
	add := func(acc, v int) int { return acc + v }

	tests := []struct {
		name     string
		input    []int
		initial  int
		expected int
	}{
		{name: "empty seq returns initial", input: []int{}, initial: 42, expected: 42},
		{name: "sum", input: []int{1, 2, 3, 4, 5}, initial: 0, expected: 15},
		{name: "non-zero initial", input: []int{1, 2, 3}, initial: 10, expected: 16},
		{name: "single element", input: []int{7}, initial: 0, expected: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			initial := tt.initial
			// when
			result := FoldSeq(slices.Values(input), initial, add)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("left order", func(t *testing.T) {
		// subtraction is non-associative: verifies left-to-right application
		// given
		input := []int{1, 2, 3}
		// when: ((10-1)-2)-3 = 4
		result := FoldSeq(slices.Values(input), 10, func(acc, v int) int { return acc - v })
		// then
		assert.Equal(t, 4, result)
	})
}
