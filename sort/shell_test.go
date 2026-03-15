package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	asc := func(a, b int) Order {
		if a < b {
			return Before
		}

		if a > b {
			return After
		}
		return Equal
	}

	dec := func(a, b int) Order {
		if a > b {
			return Before
		}

		if a < b {
			return After
		}
		return Equal
	}

	tests := []struct {
		name     string
		arr      []int
		cmd      Comparator[int]
		expected []int
	}{
		{
			name:     "sort ascending",
			arr:      []int{9, 7, 4, 3, 6, 1, 5, 2, 8, 10},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "sort descending",
			arr:      []int{9, 7, 4, 3, 6, 1, 5, 2, 8, 10},
			cmd:      dec,
			expected: []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		},
		{
			name:     "sort empty",
			arr:      []int{},
			cmd:      asc,
			expected: []int{},
		},
		{
			name:     "sort single element",
			arr:      []int{42},
			cmd:      asc,
			expected: []int{42},
		},
		{
			name:     "sort already sorted",
			arr:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "sort reverse sorted",
			arr:      []int{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		},
		{
			name:     "sort with duplicates",
			arr:      []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3},
			cmd:      asc,
			expected: []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			Shell(tt.arr, tt.cmd)
			// then
			assert.Equal(t, tt.expected, tt.arr)
		})
	}
}
