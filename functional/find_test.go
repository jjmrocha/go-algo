package functional

import (
	"slices"
	"testing"
)

func TestFind(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("nil input", func(t *testing.T) {
		got, ok := Find([]int(nil), isEven)
		if ok || got != 0 {
			t.Errorf("Find nil: got (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got, ok := Find([]int{}, isEven)
		if ok || got != 0 {
			t.Errorf("Find empty: got (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("no match", func(t *testing.T) {
		got, ok := Find([]int{1, 3, 5}, isEven)
		if ok || got != 0 {
			t.Errorf("Find no match: got (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("first element matches", func(t *testing.T) {
		got, ok := Find([]int{2, 3, 5}, isEven)
		if !ok || got != 2 {
			t.Errorf("Find first: got (%v, %v), want (2, true)", got, ok)
		}
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		got, ok := Find([]int{1, 4, 6}, isEven)
		if !ok || got != 4 {
			t.Errorf("Find first-of-many: got (%v, %v), want (4, true)", got, ok)
		}
	})

	t.Run("last element matches", func(t *testing.T) {
		got, ok := Find([]int{1, 3, 6}, isEven)
		if !ok || got != 6 {
			t.Errorf("Find last: got (%v, %v), want (6, true)", got, ok)
		}
	})
}

func TestFindSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("empty seq", func(t *testing.T) {
		got, ok := FindSeq(slices.Values([]int{}), isEven)
		if ok || got != 0 {
			t.Errorf("FindSeq empty: got (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("no match", func(t *testing.T) {
		got, ok := FindSeq(slices.Values([]int{1, 3, 5}), isEven)
		if ok || got != 0 {
			t.Errorf("FindSeq no match: got (%v, %v), want (0, false)", got, ok)
		}
	})

	t.Run("first element matches", func(t *testing.T) {
		got, ok := FindSeq(slices.Values([]int{2, 3, 5}), isEven)
		if !ok || got != 2 {
			t.Errorf("FindSeq first: got (%v, %v), want (2, true)", got, ok)
		}
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		got, ok := FindSeq(slices.Values([]int{1, 4, 6}), isEven)
		if !ok || got != 4 {
			t.Errorf("FindSeq first-of-many: got (%v, %v), want (4, true)", got, ok)
		}
	})

	t.Run("last element matches", func(t *testing.T) {
		got, ok := FindSeq(slices.Values([]int{1, 3, 6}), isEven)
		if !ok || got != 6 {
			t.Errorf("FindSeq last: got (%v, %v), want (6, true)", got, ok)
		}
	})
}
