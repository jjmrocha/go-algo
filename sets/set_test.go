package sets

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
	assert.False(t, result.Contains(1))
}

func TestOf(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{name: "nil input", input: nil, expected: []int{}},
		{name: "empty input", input: []int{}, expected: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expected: []int{1, 2, 3}},
		{name: "with duplicates", input: []int{1, 2, 1, 3, 2}, expected: []int{1, 2, 3}},
		{name: "all duplicates", input: []int{5, 5, 5}, expected: []int{5}},
		{name: "single element", input: []int{7}, expected: []int{7}},
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
	t.Run("single element", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		s.Add(1)
		// then
		assert.Equal(t, 1, s.Len())
		assert.True(t, s.Contains(1))
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		s.Add(1, 2, 3)
		// then
		assert.Equal(t, 3, s.Len())
	})

	t.Run("duplicate is idempotent", func(t *testing.T) {
		// given
		s := New[int]()
		s.Add(1)
		// when
		s.Add(1)
		// then
		assert.Equal(t, 1, s.Len())
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		s.Add()
		// then
		assert.Equal(t, 0, s.Len())
	})
}

func TestRemove(t *testing.T) {
	t.Run("existing element", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		s.Remove(2)
		// then
		assert.Equal(t, 2, s.Len())
		assert.False(t, s.Contains(2))
	})

	t.Run("absent element is a no-op", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		s.Remove(99)
		// then
		assert.Equal(t, 3, s.Len())
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3, 4})
		// when
		s.Remove(1, 3)
		// then
		assert.Equal(t, 2, s.Len())
		assert.False(t, s.Contains(1))
		assert.False(t, s.Contains(3))
		assert.True(t, s.Contains(2))
		assert.True(t, s.Contains(4))
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		s := Of([]int{1, 2})
		// when
		s.Remove()
		// then
		assert.Equal(t, 2, s.Len())
	})
}

func TestContains(t *testing.T) {
	// given
	s := Of([]int{1, 2, 3})

	tests := []struct {
		name     string
		value    int
		expected bool
	}{
		{name: "element 1 present", value: 1, expected: true},
		{name: "element 2 present", value: 2, expected: true},
		{name: "element 3 present", value: 3, expected: true},
		{name: "element 0 absent", value: 0, expected: false},
		{name: "element 4 absent", value: 4, expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			value := tt.value
			// when
			result := s.Contains(value)
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
		{name: "single", input: []int{1}, expected: 1},
		{name: "multiple", input: []int{1, 2, 3}, expected: 3},
		{name: "duplicates counted once", input: []int{1, 1, 2}, expected: 2},
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

func TestToSlice(t *testing.T) {
	t.Run("empty set returns empty slice", func(t *testing.T) {
		// when
		result := New[int]().ToSlice()
		// then
		assert.Empty(t, result)
	})

	t.Run("contains all elements", func(t *testing.T) {
		// given
		s := Of([]int{3, 1, 4, 1, 5})
		// when
		result := s.ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 3, 4, 5}, result)
	})

	t.Run("length matches Len", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := s.ToSlice()
		// then
		assert.Len(t, result, s.Len())
	})
}

// TestSet_NilSafety documents that read operations on a nil Set do not panic.
func TestSet_NilSafety(t *testing.T) {
	// given
	var s Set[int]
	// then
	assert.Equal(t, 0, s.Len())
	assert.False(t, s.Contains(1))
	result := s.ToSlice()
	assert.Empty(t, result)
	s.Remove(1) // must not panic
}

func TestString(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		// given
		var s Set[int]
		// when
		result := s.String()
		// then
		assert.Equal(t, "set(nil)", result)
	})

	t.Run("empty set", func(t *testing.T) {
		// when
		result := New[int]().String()
		// then
		assert.Equal(t, "set{}", result)
	})

	t.Run("single element", func(t *testing.T) {
		// given
		s := Of([]int{42})
		// when
		result := s.String()
		// then
		assert.Equal(t, "set{42}", result)
	})

	t.Run("multiple elements — all represented", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := s.String()
		// then
		assert.Regexp(t, `^set\{.*\}$`, result)
		assert.Contains(t, result, "1")
		assert.Contains(t, result, "2")
		assert.Contains(t, result, "3")
	})
}

func TestUnion(t *testing.T) {
	t.Run("disjoint sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{4, 5, 6})
		// when
		result := a.Union(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("overlapping sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Union(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Union(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("with empty set", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Union(New[int]()).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2}, result)
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
		b := Of([]int{3, 4})
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
	t.Run("disjoint sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{4, 5, 6})
		// when
		result := a.Intersection(b)
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("partial overlap", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Intersection(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{2, 3}, result)
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Intersection(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("with empty set", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Intersection(New[int]())
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("result is independent of s", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2})
		// when
		result := a.Intersection(b)
		a.Add(10)
		// then
		assert.False(t, result.Contains(10))
	})
}

func TestDifference(t *testing.T) {
	t.Run("disjoint sets (s-o = s)", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{4, 5, 6})
		// when
		result := a.Difference(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("partial overlap", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Difference(b).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1}, result)
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Difference(b)
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("b is superset of a", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		b := Of([]int{1, 2, 3, 4})
		// when
		result := a.Difference(b)
		// then
		assert.Equal(t, 0, result.Len())
	})

	t.Run("empty other set (s-∅ = s)", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Difference(New[int]()).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1, 2}, result)
	})

	t.Run("is not symmetric", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		ab := a.Difference(b).ToSlice()
		ba := b.Difference(a).ToSlice()
		// then
		assert.ElementsMatch(t, []int{1}, ab)
		assert.ElementsMatch(t, []int{4}, ba)
	})
}

func TestValues(t *testing.T) {
	t.Run("empty set yields nothing", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		called := false
		for range s.Values() {
			called = true
		}
		// then
		assert.False(t, called)
	})

	t.Run("all elements yielded", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := slices.Collect(s.Values())
		// then
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})

	t.Run("early termination", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3, 4, 5})
		// when
		count := 0
		for range s.Values() {
			count++
			if count == 2 {
				break
			}
		}
		// then
		assert.Equal(t, 2, count)
	})
}
