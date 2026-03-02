package functional

import "testing"

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
