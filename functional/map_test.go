package functional

import (
	"slices"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		expected []int
	}{
		{
			name:     "nil input",
			input:    nil,
			fn:       func(v int) int { return v * 2 },
			expected: []int{},
		},
		{
			name:     "empty input",
			input:    []int{},
			fn:       func(v int) int { return v * 2 },
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{3},
			fn:       func(v int) int { return v * 2 },
			expected: []int{6},
		},
		{
			name:     "multiple elements",
			input:    []int{1, 2, 3},
			fn:       func(v int) int { return v * 2 },
			expected: []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.input, tt.fn)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("Map(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	t.Run("type transform", func(t *testing.T) {
		// given: int → bool: verify T→U type change works
		input := []int{0, 1, 2}
		// when
		result := Map(input, func(v int) bool { return v > 0 })
		// then
		expected := []bool{false, true, true}
		if !slices.Equal(result, expected) {
			t.Errorf("Map type transform = %v, want %v", result, expected)
		}
	})

	t.Run("preserves order", func(t *testing.T) {
		// given
		input := []int{5, 3, 1, 4, 2}
		// when
		result := Map(input, func(v int) int { return v * 10 })
		// then
		expected := []int{50, 30, 10, 40, 20}
		if !slices.Equal(result, expected) {
			t.Errorf("Map order = %v, want %v", result, expected)
		}
	})
}

func TestMapSeq(t *testing.T) {
	double := func(v int) int { return v * 2 }

	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "empty seq", input: []int{}, expected: []int{}},
		{name: "single element", input: []int{3}, expected: []int{6}},
		{name: "multiple elements", input: []int{1, 2, 3}, expected: []int{2, 4, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(MapSeq(slices.Values(tt.input), double))
			if !slices.Equal(got, tt.expected) {
				t.Errorf("MapSeq(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}

	t.Run("type transform", func(t *testing.T) {
		// given: int → bool: verifies the output type differs from the input type
		input := []int{1, 2, 3, 4}
		// when
		result := slices.Collect(MapSeq(slices.Values(input), func(v int) bool { return v%2 == 0 }))
		// then
		expected := []bool{false, true, false, true}
		if !slices.Equal(result, expected) {
			t.Errorf("MapSeq type transform = %v, want %v", result, expected)
		}
	})
}
