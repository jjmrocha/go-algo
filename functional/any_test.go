package functional

import "testing"

func TestAny(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	tests := []struct {
		name   string
		input  []int
		expect bool
	}{
		{name: "nil input", input: nil, expect: false},
		{name: "empty input", input: []int{}, expect: false},
		{name: "no match", input: []int{1, 3, 5}, expect: false},
		{name: "first element matches", input: []int{2, 1, 3}, expect: true},
		{name: "last element matches", input: []int{1, 3, 4}, expect: true},
		{name: "all match", input: []int{2, 4, 6}, expect: true},
		{name: "single match in middle", input: []int{1, 2, 3}, expect: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Any(tt.input, isEven)
			if got != tt.expect {
				t.Errorf("Any(%v) = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}
