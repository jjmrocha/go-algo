package functional

import (
	"slices"
	"strings"
	"testing"
)

func TestZip(t *testing.T) {
	add := func(a, b int) int { return a + b }

	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{name: "nil inputs", a: nil, b: nil, expected: []int{}},
		{name: "empty inputs", a: []int{}, b: []int{}, expected: []int{}},
		{name: "equal length", a: []int{1, 2, 3}, b: []int{10, 20, 30}, expected: []int{11, 22, 33}},
		{name: "a longer than b", a: []int{1, 2, 3}, b: []int{10, 20}, expected: []int{11, 22}},
		{name: "b longer than a", a: []int{1, 2}, b: []int{10, 20, 30}, expected: []int{11, 22}},
		{name: "single element each", a: []int{5}, b: []int{3}, expected: []int{8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Zip(tt.a, tt.b, add)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("Zip(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}

	t.Run("type combine", func(t *testing.T) {
		// given: int + string → string
		a := []int{1, 2, 3}
		b := []string{"a", "b", "c"}
		// when
		result := Zip(a, b, func(n int, s string) string {
			return strings.Repeat(s, n)
		})
		// then
		expected := []string{"a", "bb", "ccc"}
		if !slices.Equal(result, expected) {
			t.Errorf("Zip type combine = %v, want %v", result, expected)
		}
	})
}
