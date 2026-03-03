package functional

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name   string
		input  []int
		expect []int
	}{
		{name: "nil input", input: nil, expect: []int{}},
		{name: "empty input", input: []int{}, expect: []int{}},
		{name: "none pass", input: []int{1, 3, 5}, expect: []int{}},
		{name: "all pass", input: []int{2, 4, 6}, expect: []int{2, 4, 6}},
		{name: "partial pass", input: []int{1, 2, 3, 4, 5}, expect: []int{2, 4}},
		{name: "single match", input: []int{1, 2, 3}, expect: []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, isEven)
			if !slices.Equal(got, tt.expect) {
				t.Errorf("Filter(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestFilter_PreservesOrder(t *testing.T) {
	input := []int{5, 1, 4, 2, 3}
	got := Filter(input, func(v int) bool { return v > 2 })
	expect := []int{5, 4, 3}
	if !slices.Equal(got, expect) {
		t.Errorf("Filter order = %v, want %v", got, expect)
	}
}

func TestFilterSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name   string
		input  []int
		expect []int
	}{
		{name: "empty", input: []int{}, expect: []int{}},
		{name: "none pass", input: []int{1, 3, 5}, expect: []int{}},
		{name: "all pass", input: []int{2, 4, 6}, expect: []int{2, 4, 6}},
		{name: "mixed", input: []int{1, 2, 3, 4, 5}, expect: []int{2, 4}},
		{name: "single match", input: []int{1, 2, 3}, expect: []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(FilterSeq(slices.Values(tt.input), isEven))
			if !slices.Equal(got, tt.expect) {
				t.Errorf("FilterSeq(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
