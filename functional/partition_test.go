package functional

import (
	"slices"
	"testing"
)

func TestPartition(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name        string
		input       []int
		expectMatch []int
		expectRest  []int
	}{
		{name: "nil input", input: nil, expectMatch: nil, expectRest: nil},
		{name: "empty input", input: []int{}, expectMatch: nil, expectRest: nil},
		{name: "all matching", input: []int{2, 4, 6}, expectMatch: []int{2, 4, 6}, expectRest: nil},
		{name: "none matching", input: []int{1, 3, 5}, expectMatch: nil, expectRest: []int{1, 3, 5}},
		{name: "mixed", input: []int{1, 2, 3, 4, 5}, expectMatch: []int{2, 4}, expectRest: []int{1, 3, 5}},
		{name: "single match", input: []int{1, 2}, expectMatch: []int{2}, expectRest: []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, rest := Partition(tt.input, isEven)
			if !slices.Equal(match, tt.expectMatch) {
				t.Errorf("Partition matching: got %v, want %v", match, tt.expectMatch)
			}
			if !slices.Equal(rest, tt.expectRest) {
				t.Errorf("Partition non-matching: got %v, want %v", rest, tt.expectRest)
			}
		})
	}
}

func TestPartition_PreservesOrder(t *testing.T) {
	input := []int{5, 2, 8, 1, 4}
	match, rest := Partition(input, func(v int) bool { return v > 3 })
	if !slices.Equal(match, []int{5, 8, 4}) {
		t.Errorf("Partition match order = %v, want [5 8 4]", match)
	}
	if !slices.Equal(rest, []int{2, 1}) {
		t.Errorf("Partition rest order = %v, want [2 1]", rest)
	}
}
