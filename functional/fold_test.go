package functional

import (
	"slices"
	"testing"
)

func TestFold(t *testing.T) {
	add := func(acc, v int) int { return acc + v }

	tests := []struct {
		name     string
		input    []int
		initial  int
		expected int
	}{
		{name: "nil input returns initial", input: nil, initial: 42, expected: 42},
		{name: "empty input returns initial", input: []int{}, initial: 7, expected: 7},
		{name: "single element", input: []int{5}, initial: 0, expected: 5},
		{name: "sum", input: []int{1, 2, 3, 4}, initial: 0, expected: 10},
		{name: "non-zero initial", input: []int{1, 2, 3}, initial: 10, expected: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fold(tt.input, tt.initial, add)
			if got != tt.expected {
				t.Errorf("Fold(%v, %d) = %d, want %d", tt.input, tt.initial, got, tt.expected)
			}
		})
	}

	t.Run("left order", func(t *testing.T) {
		// subtraction is non-commutative: verifies left-to-right application
		// given
		input := []int{1, 2, 3}
		// when: ((10 - 1) - 2) - 3 = 4
		result := Fold(input, 10, func(acc, v int) int { return acc - v })
		// then
		if result != 4 {
			t.Errorf("Fold left order = %d, want 4", result)
		}
	})

	t.Run("string concat", func(t *testing.T) {
		// given
		input := []string{"b", "c", "d"}
		// when
		result := Fold(input, "a", func(acc, v string) string { return acc + v })
		// then
		if result != "abcd" {
			t.Errorf("Fold string concat = %q, want %q", result, "abcd")
		}
	})
}

func TestFoldSeq(t *testing.T) {
	add := func(acc, v int) int { return acc + v }

	tests := []struct {
		name     string
		input    []int
		initial  int
		expected int
	}{
		{name: "empty seq returns initial", input: []int{}, initial: 42, expected: 42},
		{name: "sum", input: []int{1, 2, 3, 4, 5}, initial: 0, expected: 15},
		{name: "non-zero initial", input: []int{1, 2, 3}, initial: 10, expected: 16},
		{name: "single element", input: []int{7}, initial: 0, expected: 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FoldSeq(slices.Values(tt.input), tt.initial, add)
			if got != tt.expected {
				t.Errorf("FoldSeq(%v, %d) = %v, want %v", tt.input, tt.initial, got, tt.expected)
			}
		})
	}

	t.Run("left order", func(t *testing.T) {
		// subtraction is non-associative: verifies left-to-right application
		// given
		input := []int{1, 2, 3}
		// when: ((10-1)-2)-3 = 4
		result := FoldSeq(slices.Values(input), 10, func(acc, v int) int { return acc - v })
		// then
		if result != 4 {
			t.Errorf("FoldSeq left order: got %d, want 4", result)
		}
	})
}
