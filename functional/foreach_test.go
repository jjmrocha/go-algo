package functional

import (
	"slices"
	"testing"
)

func TestForEach(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		var got []int
		ForEach([]int{1, 2, 3}, func(v int) { got = append(got, v) })
		expect := []int{1, 2, 3}
		for i, v := range expect {
			if i >= len(got) || got[i] != v {
				t.Fatalf("ForEach order: got %v, want %v", got, expect)
			}
		}
		if len(got) != len(expect) {
			t.Fatalf("ForEach len: got %d, want %d", len(got), len(expect))
		}
	})

	t.Run("nil input — fn never called", func(t *testing.T) {
		called := false
		ForEach([]int(nil), func(_ int) { called = true })
		if called {
			t.Error("ForEach on nil input should not call fn")
		}
	})

	t.Run("empty input — fn never called", func(t *testing.T) {
		called := false
		ForEach([]int{}, func(_ int) { called = true })
		if called {
			t.Error("ForEach on empty input should not call fn")
		}
	})
}

func TestForEachSeq(t *testing.T) {
	t.Run("fn called for each element in order", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		var got []int
		ForEachSeq(slices.Values(input), func(v int) { got = append(got, v) })
		if !slices.Equal(got, input) {
			t.Errorf("ForEachSeq = %v, want %v", got, input)
		}
	})

	t.Run("empty seq — fn never called", func(t *testing.T) {
		called := false
		ForEachSeq(slices.Values([]int{}), func(_ int) { called = true })
		if called {
			t.Error("ForEachSeq called fn on empty seq")
		}
	})
}
