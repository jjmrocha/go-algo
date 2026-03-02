package functional

import "testing"

func TestAll(t *testing.T) {
	isPositive := func(v int) bool { return v > 0 }

	tests := []struct {
		name   string
		input  []int
		expect bool
	}{
		{name: "nil input — vacuously true", input: nil, expect: true},
		{name: "empty input — vacuously true", input: []int{}, expect: true},
		{name: "all match", input: []int{1, 2, 3}, expect: true},
		{name: "first element fails", input: []int{-1, 2, 3}, expect: false},
		{name: "last element fails", input: []int{1, 2, -3}, expect: false},
		{name: "none match", input: []int{-1, -2, -3}, expect: false},
		{name: "single matching element", input: []int{1}, expect: true},
		{name: "single non-matching element", input: []int{-1}, expect: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := All(tt.input, isPositive)
			if got != tt.expect {
				t.Errorf("All(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
