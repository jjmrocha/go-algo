package functional

import "testing"

func TestFold(t *testing.T) {
	add := func(acc, v int) int { return acc + v }

	tests := []struct {
		name    string
		input   []int
		initial int
		expect  int
	}{
		{name: "nil input returns initial", input: nil, initial: 42, expect: 42},
		{name: "empty input returns initial", input: []int{}, initial: 7, expect: 7},
		{name: "single element", input: []int{5}, initial: 0, expect: 5},
		{name: "sum", input: []int{1, 2, 3, 4}, initial: 0, expect: 10},
		{name: "non-zero initial", input: []int{1, 2, 3}, initial: 10, expect: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Fold(tt.input, tt.initial, add)
			if got != tt.expect {
				t.Errorf("Fold(%v, %d) = %d, want %d", tt.input, tt.initial, got, tt.expect)
			}
		})
	}
}

func TestFold_LeftOrder(t *testing.T) {
	// subtraction is non-commutative: verifies left-to-right application
	// ((10 - 1) - 2) - 3 = 4
	got := Fold([]int{1, 2, 3}, 10, func(acc, v int) int { return acc - v })
	if got != 4 {
		t.Errorf("Fold left order = %d, want 4", got)
	}
}

func TestFold_StringConcat(t *testing.T) {
	got := Fold([]string{"b", "c", "d"}, "a", func(acc, v string) string { return acc + v })
	if got != "abcd" {
		t.Errorf("Fold string concat = %q, want %q", got, "abcd")
	}
}
