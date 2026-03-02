package functional

import (
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		fn     func(int) int
		expect []int
	}{
		{
			name:   "nil input",
			input:  nil,
			fn:     func(v int) int { return v * 2 },
			expect: []int{},
		},
		{
			name:   "empty input",
			input:  []int{},
			fn:     func(v int) int { return v * 2 },
			expect: []int{},
		},
		{
			name:   "single element",
			input:  []int{3},
			fn:     func(v int) int { return v * 2 },
			expect: []int{6},
		},
		{
			name:   "multiple elements",
			input:  []int{1, 2, 3},
			fn:     func(v int) int { return v * 2 },
			expect: []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.input, tt.fn)
			if !slices.Equal(got, tt.expect) {
				t.Errorf("Map(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestMap_TypeTransform(t *testing.T) {
	// int → bool: verify T→U type change works
	got := Map([]int{0, 1, 2}, func(v int) bool { return v > 0 })
	expect := []bool{false, true, true}
	if !slices.Equal(got, expect) {
		t.Errorf("Map type transform = %v, want %v", got, expect)
	}
}

func TestMap_PreservesOrder(t *testing.T) {
	input := []int{5, 3, 1, 4, 2}
	got := Map(input, func(v int) int { return v * 10 })
	expect := []int{50, 30, 10, 40, 20}
	if !slices.Equal(got, expect) {
		t.Errorf("Map order = %v, want %v", got, expect)
	}
}
