package bag

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// when
	result := New[int]()
	// then
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Len())
	assert.True(t, result.Empty())
}

func TestOf(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "nil input", input: nil, expected: []int{}},
		{name: "empty input", input: []int{}, expected: []int{}},
		{name: "single element", input: []int{5}, expected: []int{5}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "with duplicates", input: []int{1, 2, 1}, expected: []int{1, 1, 2}},
		{name: "all same", input: []int{3, 3, 3}, expected: []int{3, 3, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := Of(input).ToSlice()
			// then
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestAdd(t *testing.T) {
	t.Run("new item", func(t *testing.T) {
		// given
		b := New[int]()
		// when
		b.Add(1)
		// then
		assert.Equal(t, 1, b.Count(1))
		assert.Equal(t, 1, b.Len())
	})

	t.Run("same item increments count", func(t *testing.T) {
		// given
		b := New[int]()
		b.Add(1)
		// when
		b.Add(1)
		// then
		assert.Equal(t, 2, b.Count(1))
		assert.Equal(t, 2, b.Len())
	})

	t.Run("multiple items at once", func(t *testing.T) {
		// given
		b := New[int]()
		// when
		b.Add(1, 2, 3)
		// then
		assert.Equal(t, 3, b.Len())
		assert.Equal(t, 1, b.Count(1))
		assert.Equal(t, 1, b.Count(2))
		assert.Equal(t, 1, b.Count(3))
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		b := New[int]()
		// when
		b.Add()
		// then
		assert.Equal(t, 0, b.Len())
	})
}

func TestRemove(t *testing.T) {
	t.Run("decrements count above one", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 1})
		// when
		b.Remove(1)
		// then
		assert.Equal(t, 2, b.Count(1))
		assert.Equal(t, 2, b.Len())
	})

	t.Run("deletes item when count reaches zero", func(t *testing.T) {
		// given
		b := Of([]int{1})
		// when
		b.Remove(1)
		// then
		assert.Equal(t, 0, b.Count(1))
		assert.False(t, b.Contains(1))
		assert.Equal(t, 0, b.Len())
	})

	t.Run("absent item is a no-op", func(t *testing.T) {
		// given
		b := Of([]int{1, 2})
		// when
		b.Remove(99)
		// then
		assert.Equal(t, 2, b.Len())
	})

	t.Run("multiple items at once", func(t *testing.T) {
		// given
		b := Of([]int{1, 2, 3})
		// when
		b.Remove(1, 3)
		// then
		assert.Equal(t, 1, b.Len())
		assert.False(t, b.Contains(1))
		assert.True(t, b.Contains(2))
		assert.False(t, b.Contains(3))
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		b := Of([]int{1, 2})
		// when
		b.Remove()
		// then
		assert.Equal(t, 2, b.Len())
	})
}

func TestClear(t *testing.T) {
	t.Run("removes all items", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 2, 3})
		// when
		b.Clear()
		// then
		assert.Equal(t, 0, b.Len())
		assert.True(t, b.Empty())
		assert.False(t, b.Contains(1))
	})

	t.Run("clearing empty bag is a no-op", func(t *testing.T) {
		// given
		b := New[int]()
		// when
		b.Clear()
		// then
		assert.Equal(t, 0, b.Len())
	})
}

func TestContains(t *testing.T) {
	// given
	b := Of([]int{1, 2, 2, 3})

	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{name: "element present once", value: 1, expected: true},
		{name: "element present multiple times", value: 2, expected: true},
		{name: "element absent", value: 99, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			value := tt.value
			// when
			result := b.Contains(value)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{name: "empty", input: []int{}, expected: 0},
		{name: "single element", input: []int{1}, expected: 1},
		{name: "multiple unique elements", input: []int{1, 2, 3}, expected: 3},
		{name: "counts duplicates", input: []int{1, 1, 2}, expected: 3},
		{name: "all same element", input: []int{5, 5, 5, 5}, expected: 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			input := tt.input
			// when
			result := Of(input).Len()
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCount(t *testing.T) {
	// given
	b := Of([]int{1, 2, 2, 2, 3})

	tests := []struct {
		name     string
		value    int
		expected int
	}{
		{name: "present once", value: 1, expected: 1},
		{name: "present multiple times", value: 2, expected: 3},
		{name: "absent", value: 99, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			value := tt.value
			// when
			result := b.Count(value)
			// then
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmpty(t *testing.T) {
	t.Run("empty bag", func(t *testing.T) {
		// when
		result := New[int]().Empty()
		// then
		assert.True(t, result)
	})

	t.Run("non-empty bag", func(t *testing.T) {
		// given
		b := Of([]int{1})
		// when
		result := b.Empty()
		// then
		assert.False(t, result)
	})
}

func TestToSlice(t *testing.T) {
	t.Run("empty bag returns empty slice", func(t *testing.T) {
		// when
		result := New[int]().ToSlice()
		// then
		assert.Empty(t, result)
	})

	t.Run("items appear with multiplicity", func(t *testing.T) {
		// given
		b := Of([]int{1, 2, 2, 3})
		// when
		result := b.ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 2, 3}, result)
	})

	t.Run("length matches Len", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 2})
		// when
		result := b.ToSlice()
		// then
		assert.Len(t, result, b.Len())
	})
}

func TestUnique(t *testing.T) {
	t.Run("empty bag returns empty slice", func(t *testing.T) {
		// when
		result := New[int]().Unique()
		// then
		assert.Empty(t, result)
	})

	t.Run("each distinct item appears once", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 2, 2, 3})
		// when
		result := b.Unique()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("length matches number of distinct items", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 2, 3, 3})
		// when
		result := b.Unique()
		// then
		assert.Len(t, result, len(b))
	})
}

func TestValues(t *testing.T) {
	t.Run("empty bag yields nothing", func(t *testing.T) {
		// given
		b := New[int]()
		// when
		called := false
		for range b.Values() {
			called = true
		}
		// then
		assert.False(t, called)
	})

	t.Run("yields all items with multiplicity", func(t *testing.T) {
		// given
		b := Of([]int{1, 2, 2, 3})
		// when
		result := slices.Collect(b.Values())
		// then
		assert.ElementsMatch(t, []int{1, 2, 2, 3}, result)
	})

	t.Run("early termination", func(t *testing.T) {
		// given
		b := Of([]int{1, 2, 3, 4, 5})
		// when
		count := 0
		for range b.Values() {
			count++
			if count == 3 {
				break
			}
		}
		// then
		assert.Equal(t, 3, count)
	})

	t.Run("multiple iterations", func(t *testing.T) {
		// given
		b := Of([]int{1, 1, 2, 3})
		first := slices.Collect(b.Values())
		// when
		result := slices.Collect(b.Values())
		// then — each range loop gets independent state
		assert.ElementsMatch(t, first, result)
	})
}

func TestString(t *testing.T) {
	t.Run("nil bag", func(t *testing.T) {
		// given
		var b Bag[int]
		// when
		result := b.String()
		// then
		assert.Equal(t, "bag(nil)", result)
	})

	t.Run("empty bag", func(t *testing.T) {
		// when
		result := New[int]().String()
		// then
		assert.Equal(t, "bag{}", result)
	})

	t.Run("single element", func(t *testing.T) {
		// given
		b := Of([]int{42})
		// when
		result := b.String()
		// then
		assert.Equal(t, "bag{42: 1}", result)
	})

	t.Run("multiple elements — all represented", func(t *testing.T) {
		// given
		b := Of([]int{1, 2, 2, 3})
		// when
		result := b.String()
		// then
		assert.Regexp(t, `^bag\{.*\}$`, result)
		assert.Contains(t, result, "1: 1")
		assert.Contains(t, result, "2: 2")
		assert.Contains(t, result, "3: 1")
	})
}

func TestUnion(t *testing.T) {
	t.Run("disjoint bags preserves all counts", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 2})
		b := Of([]int{3, 4})
		// when
		result := a.Union(b)
		// then
		assert.Equal(t, 1, result.Count(1))
		assert.Equal(t, 2, result.Count(2))
		assert.Equal(t, 1, result.Count(3))
		assert.Equal(t, 1, result.Count(4))
		assert.Equal(t, 5, result.Len())
	})

	t.Run("overlapping bags sums counts", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 2})
		b := Of([]int{2, 3})
		// when
		result := a.Union(b)
		// then
		assert.Equal(t, 1, result.Count(1))
		assert.Equal(t, 3, result.Count(2))
		assert.Equal(t, 1, result.Count(3))
	})

	t.Run("with empty bag", func(t *testing.T) {
		// given
		a := Of([]int{1, 1, 2})
		// when
		result := a.Union(New[int]())
		// then
		assert.Equal(t, 2, result.Count(1))
		assert.Equal(t, 1, result.Count(2))
		assert.Equal(t, 3, result.Len())
	})

	t.Run("both empty", func(t *testing.T) {
		// when
		result := New[int]().Union(New[int]())
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("result is independent of operands", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		b := Of([]int{3})
		// when
		result := a.Union(b)
		a.Add(10)
		b.Add(20)
		// then
		assert.False(t, result.Contains(10))
		assert.False(t, result.Contains(20))
	})
}

func TestIntersection(t *testing.T) {
	t.Run("disjoint bags yields empty", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		b := Of([]int{3, 4})
		// when
		result := a.Intersection(b)
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("overlapping bags takes minimum count", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 2, 3})
		b := Of([]int{2, 2, 2, 3, 3})
		// when
		result := a.Intersection(b)
		// then
		assert.Equal(t, 2, result.Count(2))
		assert.Equal(t, 1, result.Count(3))
		assert.False(t, result.Contains(1))
		assert.Equal(t, 3, result.Len())
	})

	t.Run("identical bags returns same counts", func(t *testing.T) {
		// given
		a := Of([]int{1, 1, 2})
		b := Of([]int{1, 1, 2})
		// when
		result := a.Intersection(b)
		// then
		assert.Equal(t, 2, result.Count(1))
		assert.Equal(t, 1, result.Count(2))
	})

	t.Run("with empty bag yields empty", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Intersection(New[int]())
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("result is independent of operands", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		b := Of([]int{1, 2})
		// when
		result := a.Intersection(b)
		a.Add(10)
		// then
		assert.False(t, result.Contains(10))
	})
}

// TestBag_NilSafety documents that read operations on a nil Bag do not panic.
func TestBag_NilSafety(t *testing.T) {
	// given
	var b Bag[int]
	// then
	assert.Equal(t, 0, b.Len())
	assert.Equal(t, 0, b.Count(1))
	assert.False(t, b.Contains(1))
	assert.True(t, b.Empty())
	b.Remove(1) // must not panic
}
