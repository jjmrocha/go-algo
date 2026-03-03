package sets

import (
	"slices"
	"testing"
)

func TestNew(t *testing.T) {
	s := New[int]()
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.Len() != 0 {
		t.Errorf("New: Len = %d, want 0", s.Len())
	}
	if s.Contains(1) {
		t.Error("New: Contains(1) = true, want false")
	}
}

func TestFromSlice(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect []int // sorted for comparison
	}{
		{name: "nil input", input: nil, expect: []int{}},
		{name: "empty input", input: []int{}, expect: []int{}},
		{name: "no duplicates", input: []int{1, 2, 3}, expect: []int{1, 2, 3}},
		{name: "with duplicates", input: []int{1, 2, 1, 3, 2}, expect: []int{1, 2, 3}},
		{name: "all duplicates", input: []int{5, 5, 5}, expect: []int{5}},
		{name: "single element", input: []int{7}, expect: []int{7}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlice(tt.input).ToSlice()
			slices.Sort(got)
			if !slices.Equal(got, tt.expect) {
				t.Errorf("FromSlice(%v) sorted = %v, want %v", tt.input, got, tt.expect)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Run("single element", func(t *testing.T) {
		s := New[int]()
		s.Add(1)
		if s.Len() != 1 {
			t.Errorf("Len = %d, want 1", s.Len())
		}
		if !s.Contains(1) {
			t.Error("Contains(1) = false, want true")
		}
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		s := New[int]()
		s.Add(1, 2, 3)
		if s.Len() != 3 {
			t.Errorf("Len = %d, want 3", s.Len())
		}
	})

	t.Run("duplicate is idempotent", func(t *testing.T) {
		s := New[int]()
		s.Add(1)
		s.Add(1)
		if s.Len() != 1 {
			t.Errorf("Len = %d, want 1 after duplicate Add", s.Len())
		}
	})

	t.Run("no items is a no-op", func(t *testing.T) {
		s := New[int]()
		s.Add()
		if s.Len() != 0 {
			t.Errorf("Len = %d, want 0", s.Len())
		}
	})
}

func TestRemove(t *testing.T) {
	t.Run("existing element", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		s.Remove(2)
		if s.Len() != 2 {
			t.Errorf("Len = %d, want 2", s.Len())
		}
		if s.Contains(2) {
			t.Error("Contains(2) = true after Remove, want false")
		}
	})

	t.Run("absent element is a no-op", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		s.Remove(99)
		if s.Len() != 3 {
			t.Errorf("Len = %d, want 3", s.Len())
		}
	})

	t.Run("multiple elements at once", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3, 4})
		s.Remove(1, 3)
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
		s := FromSlice([]int{1, 2})
		s.Remove()
		if s.Len() != 2 {
			t.Errorf("Len = %d, want 2", s.Len())
		}
	})
}

func TestContains(t *testing.T) {
	s := FromSlice([]int{1, 2, 3})

	tests := []struct {
		value  int
		expect bool
	}{
		{value: 1, expect: true},
		{value: 2, expect: true},
		{value: 3, expect: true},
		{value: 0, expect: false},
		{value: 4, expect: false},
	}

	for _, tt := range tests {
		got := s.Contains(tt.value)
		if got != tt.expect {
			t.Errorf("Contains(%d) = %v, want %v", tt.value, got, tt.expect)
		}
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		expect int
	}{
		{name: "empty", input: []int{}, expect: 0},
		{name: "single", input: []int{1}, expect: 1},
		{name: "multiple", input: []int{1, 2, 3}, expect: 3},
		{name: "duplicates counted once", input: []int{1, 1, 2}, expect: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromSlice(tt.input).Len()
			if got != tt.expect {
				t.Errorf("Len = %d, want %d", got, tt.expect)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	t.Run("empty set returns empty slice", func(t *testing.T) {
		got := New[int]().ToSlice()
		if len(got) != 0 {
			t.Errorf("ToSlice on empty set = %v, want empty", got)
		}
	})

	t.Run("contains all elements", func(t *testing.T) {
		s := FromSlice([]int{3, 1, 4, 1, 5})
		got := s.ToSlice()
		slices.Sort(got)
		expect := []int{1, 3, 4, 5}
		if !slices.Equal(got, expect) {
			t.Errorf("ToSlice sorted = %v, want %v", got, expect)
		}
	})

	t.Run("length matches Len", func(t *testing.T) {
		s := FromSlice([]int{1, 2, 3})
		got := s.ToSlice()
		if len(got) != s.Len() {
			t.Errorf("len(ToSlice) = %d, want %d", len(got), s.Len())
		}
	})
}

// TestSet_NilSafety documents that read operations on a nil Set do not panic.
func TestSet_NilSafety(t *testing.T) {
	var s Set[int]

	if s.Len() != 0 {
		t.Errorf("nil.Len() = %d, want 0", s.Len())
	}
	if s.Contains(1) {
		t.Error("nil.Contains(1) = true, want false")
	}
	got := s.ToSlice()
	if len(got) != 0 {
		t.Errorf("nil.ToSlice() = %v, want empty", got)
	}
	s.Remove(1) // must not panic
}
