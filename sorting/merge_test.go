package sorting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	asc := func(a, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	}

	desc := func(a, b int) int {
		if a > b {
			return -1
		}
		if a < b {
			return 1
		}
		return 0
	}

	tests := []struct {
		name     string
		input    []int
		cmp      Comparator[int]
		expected []int
	}{
		{
			name:     "sort ascending",
			input:    []int{4, 3, 6, 1, 5, 2},
			cmp:      asc,
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "sort descending",
			input:    []int{4, 3, 6, 1, 5, 2},
			cmp:      desc,
			expected: []int{6, 5, 4, 3, 2, 1},
		},
		{
			name:     "empty slice",
			input:    []int{},
			cmp:      asc,
			expected: []int{},
		},
		{
			name:     "single element",
			input:    []int{42},
			cmp:      asc,
			expected: []int{42},
		},
		{
			name:     "already sorted",
			input:    []int{1, 2, 3, 4, 5},
			cmp:      asc,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reverse sorted",
			input:    []int{5, 4, 3, 2, 1},
			cmp:      asc,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "duplicates",
			input:    []int{1, 5, 4, 5, 2, 2, 4, 1},
			cmp:      asc,
			expected: []int{1, 1, 2, 2, 4, 4, 5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			Merge(input, tt.cmp)
			// then
			assert.Equal(t, tt.expected, input)
		})
	}
}
