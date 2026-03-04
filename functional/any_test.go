package functional

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAny(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{name: "nil input", input: nil, expected: false},
		{name: "empty input", input: []int{}, expected: false},
		{name: "no match", input: []int{1, 3, 5}, expected: false},
		{name: "first element matches", input: []int{2, 1, 3}, expected: true},
		{name: "last element matches", input: []int{1, 3, 4}, expected: true},
		{name: "all match", input: []int{2, 4, 6}, expected: true},
		{name: "single match in middle", input: []int{1, 2, 3}, expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := Any(input, isEven)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnySeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{name: "empty seq", input: []int{}, expected: false},
		{name: "no match", input: []int{1, 3, 5}, expected: false},
		{name: "first matches", input: []int{2, 3, 5}, expected: true},
		{name: "last matches", input: []int{1, 3, 4}, expected: true},
		{name: "all match", input: []int{2, 4, 6}, expected: true},
		{name: "single match", input: []int{2}, expected: true},
		{name: "single no match", input: []int{1}, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := AnySeq(slices.Values(input), isEven)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}
