package sets

import (
	"slices"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// when
	result := New[int]()
	// then
	if result == nil {
		t.Fatal("New returned nil")
	}
	if result.Len() != 0 {
		t.Errorf("New: Len = %d, want 0", result.Len())
	}
	if result.Contains(1) {
		t.Error("New: Contains(1) = true, want false")
	}
}

func TestOf(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int // sorted for comparison
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
			got := Of(tt.input).ToSlice()
			slices.Sort(got)
			if !slices.Equal(got, tt.expected) {
				t.Errorf("Of(%v) sorted = %v, want %v", tt.input, got, tt.expected)
			}
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
		if s.Len() != 1 {
			t.Errorf("Len = %d, want 1", s.Len())
		}
		if !s.Contains(1) {
			t.Error("Contains(1) = false, want true")
		}
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		s.Add(1, 2, 3)
		// then
		if s.Len() != 3 {
			t.Errorf("Len = %d, want 3", s.Len())
		}
	})

	t.Run("duplicate is idempotent", func(t *testing.T) {
		// given
		s := New[int]()
		s.Add(1)
		// when
		s.Add(1)
		// then
		if s.Len() != 1 {
			t.Errorf("Len = %d, want 1 after duplicate Add", s.Len())
		}
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		s := New[int]()
		// when
		s.Add()
		// then
		if s.Len() != 0 {
			t.Errorf("Len = %d, want 0", s.Len())
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("existing element", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		s.Remove(2)
		// then
		if s.Len() != 2 {
			t.Errorf("Len = %d, want 2", s.Len())
		}
		if s.Contains(2) {
			t.Error("Contains(2) = true after Remove, want false")
		}
	})

	t.Run("absent element is a no-op", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		s.Remove(99)
		// then
		if s.Len() != 3 {
			t.Errorf("Len = %d, want 3", s.Len())
		}
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3, 4})
		// when
		s.Remove(1, 3)
		// then
		if s.Len() != 2 {
			t.Errorf("Len = %d, want 2", s.Len())
		}
		if s.Contains(1) || s.Contains(3) {
			t.Error("removed elements still present in set")
		}
		if !s.Contains(2) || !s.Contains(4) {
			t.Error("non-removed elements missing from set")
		}
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		// given
		s := Of([]int{1, 2})
		// when
		s.Remove()
		// then
		if s.Len() != 2 {
			t.Errorf("Len = %d, want 2", s.Len())
		}
	})
}

func TestContains(t *testing.T) {
	s := Of([]int{1, 2, 3})

	tests := []struct {
		value    int
		expected bool
	}{
		{value: 1, expected: true},
		{value: 2, expected: true},
		{value: 3, expected: true},
		{value: 0, expected: false},
		{value: 4, expected: false},
	}

	for _, tt := range tests {
		got := s.Contains(tt.value)
		if got != tt.expected {
			t.Errorf("Contains(%d) = %v, want %v", tt.value, got, tt.expected)
		}
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
			got := Of(tt.input).Len()
			if got != tt.expected {
				t.Errorf("Len = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	t.Run("empty set returns empty slice", func(t *testing.T) {
		// when
		result := New[int]().ToSlice()
		// then
		if len(result) != 0 {
			t.Errorf("ToSlice on empty set = %v, want empty", result)
		}
	})

	t.Run("contains all elements", func(t *testing.T) {
		// given
		s := Of([]int{3, 1, 4, 1, 5})
		// when
		result := s.ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 3, 4, 5}
		if !slices.Equal(result, expected) {
			t.Errorf("ToSlice sorted = %v, want %v", result, expected)
		}
	})

	t.Run("length matches Len", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := s.ToSlice()
		// then
		if len(result) != s.Len() {
			t.Errorf("len(ToSlice) = %d, want %d", len(result), s.Len())
		}
	})
}

// TestSet_NilSafety documents that read operations on a nil Set do not panic.
func TestSet_NilSafety(t *testing.T) {
	// given
	var s Set[int]
	// then
	if s.Len() != 0 {
		t.Errorf("nil.Len() = %d, want 0", s.Len())
	}
	if s.Contains(1) {
		t.Error("nil.Contains(1) = true, want false")
	}
	result := s.ToSlice()
	if len(result) != 0 {
		t.Errorf("nil.ToSlice() = %v, want empty", result)
	}
	s.Remove(1) // must not panic
}

func TestString(t *testing.T) {
	t.Run("nil set", func(t *testing.T) {
		// given
		var s Set[int]
		// when
		result := s.String()
		// then
		if result != "set(nil)" {
			t.Errorf("nil.String() = %q, want %q", result, "set(nil)")
		}
	})

	t.Run("empty set", func(t *testing.T) {
		// when
		result := New[int]().String()
		// then
		if result != "set{}" {
			t.Errorf("empty.String() = %q, want %q", result, "set{}")
		}
	})

	t.Run("single element", func(t *testing.T) {
		// given
		s := Of([]int{42})
		// when
		result := s.String()
		// then
		if result != "set{42}" {
			t.Errorf("single.String() = %q, want %q", result, "set{42}")
		}
	})

	t.Run("multiple elements — all represented", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := s.String()
		// then
		if !strings.HasPrefix(result, "set{") || !strings.HasSuffix(result, "}") {
			t.Errorf("String() = %q, unexpected format", result)
		}
		for _, elem := range []string{"1", "2", "3"} {
			if !strings.Contains(result, elem) {
				t.Errorf("String() = %q, missing element %s", result, elem)
			}
		}
	})
}

func TestUnion(t *testing.T) {
	t.Run("disjoint sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{4, 5, 6})
		// when
		result := a.Union(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3, 4, 5, 6}
		if !slices.Equal(result, expected) {
			t.Errorf("Union disjoint = %v, want %v", result, expected)
		}
	})

	t.Run("overlapping sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Union(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3, 4}
		if !slices.Equal(result, expected) {
			t.Errorf("Union overlapping = %v, want %v", result, expected)
		}
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Union(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Union identical = %v, want %v", result, expected)
		}
	})

	t.Run("with empty set", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Union(New[int]()).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2}
		if !slices.Equal(result, expected) {
			t.Errorf("Union with empty = %v, want %v", result, expected)
		}
	})

	t.Run("both empty", func(t *testing.T) {
		// when
		result := New[int]().Union(New[int]())
		// then
		if result.Len() != 0 {
			t.Error("Union of two empty sets should be empty")
		}
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
		if result.Contains(10) || result.Contains(20) {
			t.Error("Union result shares backing store with an operand")
		}
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
		if result.Len() != 0 {
			t.Error("Intersection of disjoint sets should be empty")
		}
	})

	t.Run("partial overlap", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Intersection(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{2, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Intersection partial = %v, want %v", result, expected)
		}
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Intersection(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Intersection identical = %v, want %v", result, expected)
		}
	})

	t.Run("with empty set", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Intersection(New[int]())
		// then
		if result.Len() != 0 {
			t.Error("Intersection with empty set should be empty")
		}
	})

	t.Run("result is independent of s", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2})
		// when
		result := a.Intersection(b)
		a.Add(10)
		// then
		if result.Contains(10) {
			t.Error("Intersection result shares backing store with s")
		}
	})
}

func TestDifference(t *testing.T) {
	t.Run("disjoint sets (s-o = s)", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{4, 5, 6})
		// when
		result := a.Difference(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Difference disjoint = %v, want %v", result, expected)
		}
	})

	t.Run("partial overlap", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		result := a.Difference(b).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1}
		if !slices.Equal(result, expected) {
			t.Errorf("Difference partial = %v, want %v", result, expected)
		}
	})

	t.Run("identical sets", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{1, 2, 3})
		// when
		result := a.Difference(b)
		// then
		if result.Len() != 0 {
			t.Error("Difference of identical sets should be empty")
		}
	})

	t.Run("b is superset of a", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		b := Of([]int{1, 2, 3, 4})
		// when
		result := a.Difference(b)
		// then
		if result.Len() != 0 {
			t.Error("Difference when b is superset of a should be empty")
		}
	})

	t.Run("empty other set (s-∅ = s)", func(t *testing.T) {
		// given
		a := Of([]int{1, 2})
		// when
		result := a.Difference(New[int]()).ToSlice()
		slices.Sort(result)
		// then
		expected := []int{1, 2}
		if !slices.Equal(result, expected) {
			t.Errorf("Difference with empty b = %v, want %v", result, expected)
		}
	})

	t.Run("is not symmetric", func(t *testing.T) {
		// given
		a := Of([]int{1, 2, 3})
		b := Of([]int{2, 3, 4})
		// when
		ab := a.Difference(b).ToSlice()
		ba := b.Difference(a).ToSlice()
		slices.Sort(ab)
		slices.Sort(ba)
		// then
		if !slices.Equal(ab, []int{1}) {
			t.Errorf("a∖b = %v, want [1]", ab)
		}
		if !slices.Equal(ba, []int{4}) {
			t.Errorf("b∖a = %v, want [4]", ba)
		}
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
		if called {
			t.Error("Values on empty set should yield nothing")
		}
	})

	t.Run("all elements yielded", func(t *testing.T) {
		// given
		s := Of([]int{1, 2, 3})
		// when
		result := slices.Collect(s.Values())
		slices.Sort(result)
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(result, expected) {
			t.Errorf("Values = %v, want %v", result, expected)
		}
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
		if count != 2 {
			t.Errorf("early termination: got %d iterations, want 2", count)
		}
	})
}
