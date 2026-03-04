package functional

import (
	"slices"
	"testing"
)

func TestFind(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("nil input", func(t *testing.T) {
		// when
		result, ok := Find([]int(nil), isEven)
		// then
		if ok || result != 0 {
			t.Errorf("Find nil: got (%v, %v), want (0, false)", result, ok)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		// when
		result, ok := Find([]int{}, isEven)
		// then
		if ok || result != 0 {
			t.Errorf("Find empty: got (%v, %v), want (0, false)", result, ok)
		}
	})

	t.Run("no match", func(t *testing.T) {
		// given
		input := []int{1, 3, 5}
		// when
		result, ok := Find(input, isEven)
		// then
		if ok || result != 0 {
			t.Errorf("Find no match: got (%v, %v), want (0, false)", result, ok)
		}
	})

	t.Run("first element matches", func(t *testing.T) {
		// given
		input := []int{2, 3, 5}
		// when
		result, ok := Find(input, isEven)
		// then
		if !ok || result != 2 {
			t.Errorf("Find first: got (%v, %v), want (2, true)", result, ok)
		}
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		// given
		input := []int{1, 4, 6}
		// when
		result, ok := Find(input, isEven)
		// then
		if !ok || result != 4 {
			t.Errorf("Find first-of-many: got (%v, %v), want (4, true)", result, ok)
		}
	})

	t.Run("last element matches", func(t *testing.T) {
		// given
		input := []int{1, 3, 6}
		// when
		result, ok := Find(input, isEven)
		// then
		if !ok || result != 6 {
			t.Errorf("Find last: got (%v, %v), want (6, true)", result, ok)
		}
	})
}

func TestFindSeq(t *testing.T) {
	isEven := func(v int) bool { return v%2 == 0 }

	t.Run("empty seq", func(t *testing.T) {
		// when
		result, ok := FindSeq(slices.Values([]int{}), isEven)
		// then
		if ok || result != 0 {
			t.Errorf("FindSeq empty: got (%v, %v), want (0, false)", result, ok)
		}
	})

	t.Run("no match", func(t *testing.T) {
		// given
		input := []int{1, 3, 5}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		if ok || result != 0 {
			t.Errorf("FindSeq no match: got (%v, %v), want (0, false)", result, ok)
		}
	})

	t.Run("first element matches", func(t *testing.T) {
		// given
		input := []int{2, 3, 5}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		if !ok || result != 2 {
			t.Errorf("FindSeq first: got (%v, %v), want (2, true)", result, ok)
		}
	})

	t.Run("returns first of multiple matches", func(t *testing.T) {
		// given
		input := []int{1, 4, 6}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		if !ok || result != 4 {
			t.Errorf("FindSeq first-of-many: got (%v, %v), want (4, true)", result, ok)
		}
	})

	t.Run("last element matches", func(t *testing.T) {
		// given
		input := []int{1, 3, 6}
		// when
		result, ok := FindSeq(slices.Values(input), isEven)
		// then
		if !ok || result != 6 {
			t.Errorf("FindSeq last: got (%v, %v), want (6, true)", result, ok)
		}
	})
}
