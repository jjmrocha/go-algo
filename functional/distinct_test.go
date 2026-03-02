package functional

import (
	"slices"
	"testing"
)

func TestDistinct(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect []int
	}{
		{name: "nil input", input: nil, expect: []int{}},
		{name: "empty input", input: []int{}, expect: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expect: []int{1, 2, 3}},
		{name: "all duplicates", input: []int{7, 7, 7}, expect: []int{7}},
		{name: "mixed", input: []int{1, 2, 1, 3, 2, 4}, expect: []int{1, 2, 3, 4}},
		{name: "single element", input: []int{5}, expect: []int{5}},
		{name: "duplicates at end", input: []int{1, 2, 3, 3}, expect: []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Distinct(tt.input)
			if !slices.Equal(got, tt.expect) {
				t.Errorf("Distinct(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestDistinct_PreservesFirstOccurrence(t *testing.T) {
	// order of first occurrences must be preserved
	got := Distinct([]int{3, 1, 4, 1, 5, 9, 2, 6, 5})
	expect := []int{3, 1, 4, 5, 9, 2, 6}
	if !slices.Equal(got, expect) {
		t.Errorf("Distinct order = %v, want %v", got, expect)
	}
}
