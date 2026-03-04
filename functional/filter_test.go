package functional

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "nil input", input: nil, expected: []int{}},
		{name: "empty input", input: []int{}, expected: []int{}},
		{name: "none pass", input: []int{1, 3, 5}, expected: []int{}},
		{name: "all pass", input: []int{2, 4, 6}, expected: []int{2, 4, 6}},
		{name: "partial pass", input: []int{1, 2, 3, 4, 5}, expected: []int{2, 4}},
		{name: "single match", input: []int{1, 2, 3}, expected: []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, isEven)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("Filter(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	t.Run("preserves order", func(t *testing.T) {
		// given
		input := []int{5, 1, 4, 2, 3}
		// when
		result := Filter(input, func(v int) bool { return v > 2 })
		// then
		expected := []int{5, 4, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Filter order = %v, want %v", result, expected)
		}
	})
}

func TestFilterSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty", input: []int{}, expected: []int{}},
		{name: "none pass", input: []int{1, 3, 5}, expected: []int{}},
		{name: "all pass", input: []int{2, 4, 6}, expected: []int{2, 4, 6}},
		{name: "mixed", input: []int{1, 2, 3, 4, 5}, expected: []int{2, 4}},
		{name: "single match", input: []int{1, 2, 3}, expected: []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(FilterSeq(slices.Values(tt.input), isEven))
			if !slices.Equal(got, tt.expected) {
				t.Errorf("FilterSeq(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
