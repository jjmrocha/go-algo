package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuick(t *testing.T) {
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
			arr:      []int{4, 3, 6, 1, 5, 2},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "sort descending",
			arr:      []int{4, 3, 6, 1, 5, 2},
			cmd:      dec,
			expected: []int{6, 5, 4, 3, 2, 1},
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
			arr:      []int{1, 2, 3, 4, 5},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "sort reverse sorted",
			arr:      []int{5, 4, 3, 2, 1},
			cmd:      asc,
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "sort with duplicates",
			arr:      []int{1, 5, 4, 5, 2, 2, 4, 1},
			cmd:      asc,
			expected: []int{1, 1, 2, 2, 4, 4, 5, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			Quick(tt.arr, tt.cmd)
			// then
			assert.Equal(t, tt.expected, tt.arr)
		})
	}
}
