package functional

import (
	"slices"
	"testing"
)

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

func TestAllSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name   string
		input  []int
		expect bool
	}{
		{name: "empty seq — vacuously true", input: []int{}, expect: true},
		{name: "all match", input: []int{2, 4, 6}, expect: true},
		{name: "first element fails", input: []int{1, 2, 4}, expect: false},
		{name: "last element fails", input: []int{2, 4, 5}, expect: false},
		{name: "none match", input: []int{1, 3, 5}, expect: false},
		{name: "single matching", input: []int{2}, expect: true},
		{name: "single non-matching", input: []int{1}, expect: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AllSeq(slices.Values(tt.input), isEven)
			if got != tt.expect {
				t.Errorf("AllSeq(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
