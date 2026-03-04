package functional

import (
	"slices"
	"testing"
)

func TestForEach(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		// given
		var collected []int
		// when
		ForEach([]int{1, 2, 3}, func(v int) { collected = append(collected, v) })
		// then
		expected := []int{1, 2, 3}
		if !slices.Equal(collected, expected) {
			t.Fatalf("ForEach order: got %v, want %v", collected, expected)
		}
	})

	t.Run("nil input — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEach([]int(nil), func(_ int) { called = true })
		// then
		if called {
			t.Error("ForEach on nil input should not call fn")
		}
	})

	t.Run("empty input — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEach([]int{}, func(_ int) { called = true })
		// then
		if called {
			t.Error("ForEach on empty input should not call fn")
		}
	})
}

func TestForEachSeq(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		// given
		input := []int{1, 2, 3, 4, 5}
		var collected []int
		// when
		ForEachSeq(slices.Values(input), func(v int) { collected = append(collected, v) })
		// then
		if !slices.Equal(collected, input) {
			t.Errorf("ForEachSeq = %v, want %v", collected, input)
		}
	})

	t.Run("empty seq — fn never called", func(t *testing.T) {
		// given
		called := false
		// when
		ForEachSeq(slices.Values([]int{}), func(_ int) { called = true })
		// then
		if called {
			t.Error("ForEachSeq called fn on empty seq")
		}
	})
}
