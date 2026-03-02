package functional

import (
	"slices"
	"testing"
)

func TestZip(t *testing.T) {
	add := func(a, b int) int { return a + b }

	tests := []struct {
		name   string
		a      []int
		b      []int
		expect []int
	}{
		{name: "nil inputs", a: nil, b: nil, expect: []int{}},
		{name: "empty inputs", a: []int{}, b: []int{}, expect: []int{}},
		{name: "equal length", a: []int{1, 2, 3}, b: []int{10, 20, 30}, expect: []int{11, 22, 33}},
		{name: "a longer than b", a: []int{1, 2, 3}, b: []int{10, 20}, expect: []int{11, 22}},
		{name: "b longer than a", a: []int{1, 2}, b: []int{10, 20, 30}, expect: []int{11, 22}},
		{name: "single element each", a: []int{5}, b: []int{3}, expect: []int{8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Zip(tt.a, tt.b, add)
			if !slices.Equal(got, tt.expect) {
				t.Errorf("Zip(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expect)
			}
		})
	}
}

func TestZip_TypeCombine(t *testing.T) {
	// int + string → string
	got := Zip([]int{1, 2, 3}, []string{"a", "b", "c"}, func(n int, s string) string {
		result := ""
		for range n {
			result += s
		}
		return result
	})
	expect := []string{"a", "bb", "ccc"}
	if !slices.Equal(got, expect) {
		t.Errorf("Zip type combine = %v, want %v", got, expect)
	}
}
