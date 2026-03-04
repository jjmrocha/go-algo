package functional

import (
	"slices"
	"testing"
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
			got := Distinct(tt.input)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("Distinct(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	t.Run("preserves first occurrence", func(t *testing.T) {
		// given
		input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
		// when
		result := Distinct(input)
		// then
		expected := []int{3, 1, 4, 5, 9, 2, 6}
		if !slices.Equal(result, expected) {
			t.Errorf("Distinct order = %v, want %v", result, expected)
		}
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
			got := slices.Collect(DistinctSeq(slices.Values(tt.input)))
			if !slices.Equal(got, tt.expected) {
				t.Errorf("DistinctSeq(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	// TestDistinctSeq_MultipleIterations verifies that each iteration starts with
	// a fresh seen set — state must not leak across iterations.
	t.Run("multiple iterations", func(t *testing.T) {
		// given
		input := []int{1, 2, 1, 3}
		seq := DistinctSeq(slices.Values(input))
		// when
		first := slices.Collect(seq)
		second := slices.Collect(seq)
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(first, expected) {
			t.Errorf("first iteration = %v, want %v", first, expected)
		}
		if !slices.Equal(second, expected) {
			t.Errorf("second iteration = %v, want %v (state leaked across iterations)", second, expected)
		}
	})
}
