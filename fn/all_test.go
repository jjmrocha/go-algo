package fn

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	isPositive := func(v int) bool { return v > 0 }

	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{name: "nil input — vacuously true", input: nil, expected: true},
		{name: "empty input — vacuously true", input: []int{}, expected: true},
		{name: "all match", input: []int{1, 2, 3}, expected: true},
		{name: "first element fails", input: []int{-1, 2, 3}, expected: false},
		{name: "last element fails", input: []int{1, 2, -3}, expected: false},
		{name: "none match", input: []int{-1, -2, -3}, expected: false},
		{name: "single matching element", input: []int{1}, expected: true},
		{name: "single non-matching element", input: []int{-1}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := All(input, isPositive)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{name: "empty seq — vacuously true", input: []int{}, expected: true},
		{name: "all match", input: []int{2, 4, 6}, expected: true},
		{name: "first element fails", input: []int{1, 2, 4}, expected: false},
		{name: "last element fails", input: []int{2, 4, 5}, expected: false},
		{name: "none match", input: []int{1, 3, 5}, expected: false},
		{name: "single matching", input: []int{2}, expected: true},
		{name: "single non-matching", input: []int{1}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := AllSeq(slices.Values(input), isEven)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
