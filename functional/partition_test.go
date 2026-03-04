package functional

import (
	"slices"
	"testing"
)

func TestPartition(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name          string
		input         []int
		expectedMatch []int
		expectedRest  []int
	}{
		{name: "nil input", input: nil, expectedMatch: nil, expectedRest: nil},
		{name: "empty input", input: []int{}, expectedMatch: nil, expectedRest: nil},
		{name: "all matching", input: []int{2, 4, 6}, expectedMatch: []int{2, 4, 6}, expectedRest: nil},
		{name: "none matching", input: []int{1, 3, 5}, expectedMatch: nil, expectedRest: []int{1, 3, 5}},
		{name: "mixed", input: []int{1, 2, 3, 4, 5}, expectedMatch: []int{2, 4}, expectedRest: []int{1, 3, 5}},
		{name: "single match", input: []int{1, 2}, expectedMatch: []int{2}, expectedRest: []int{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, rest := Partition(tt.input, isEven)
			if !slices.Equal(match, tt.expectedMatch) {
				t.Errorf("Partition matching: got %v, want %v", match, tt.expectedMatch)
			}
			if !slices.Equal(rest, tt.expectedRest) {
				t.Errorf("Partition non-matching: got %v, want %v", rest, tt.expectedRest)
			}
		})
	}

	t.Run("preserves order", func(t *testing.T) {
		// given
		input := []int{5, 2, 8, 1, 4}
		// when
		match, rest := Partition(input, func(v int) bool { return v > 3 })
		// then
		if !slices.Equal(match, []int{5, 8, 4}) {
			t.Errorf("Partition match order = %v, want [5 8 4]", match)
		}
		if !slices.Equal(rest, []int{2, 1}) {
			t.Errorf("Partition rest order = %v, want [2 1]", rest)
		}
	})
}
