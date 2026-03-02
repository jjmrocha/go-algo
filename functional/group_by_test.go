package functional

import (
	"slices"
	"testing"
)

func TestGroupBy(t *testing.T) {
	parity := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}

	t.Run("nil input", func(t *testing.T) {
		got := GroupBy([]int(nil), parity)
		if len(got) != 0 {
			t.Errorf("GroupBy nil: got %v, want empty map", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got := GroupBy([]int{}, parity)
		if len(got) != 0 {
			t.Errorf("GroupBy empty: got %v, want empty map", got)
		}
	})

	t.Run("all same key", func(t *testing.T) {
		got := GroupBy([]int{2, 4, 6}, parity)
		if len(got) != 1 {
			t.Fatalf("expected 1 group, got %d", len(got))
		}
		if !slices.Equal(got["even"], []int{2, 4, 6}) {
			t.Errorf("GroupBy all-same: got %v, want [2 4 6]", got["even"])
		}
	})

	t.Run("all different keys", func(t *testing.T) {
		got := GroupBy([]int{1, 2, 3}, func(v int) int { return v })
		if len(got) != 3 {
			t.Fatalf("expected 3 groups, got %d", len(got))
		}
		for _, v := range []int{1, 2, 3} {
			if !slices.Equal(got[v], []int{v}) {
				t.Errorf("GroupBy key %d: got %v, want [%d]", v, got[v], v)
			}
		}
	})

	t.Run("mixed — order within group preserved", func(t *testing.T) {
		got := GroupBy([]int{1, 3, 2, 5, 4}, parity)
		if !slices.Equal(got["odd"], []int{1, 3, 5}) {
			t.Errorf("GroupBy odd order: got %v, want [1 3 5]", got["odd"])
		}
		if !slices.Equal(got["even"], []int{2, 4}) {
			t.Errorf("GroupBy even order: got %v, want [2 4]", got["even"])
		}
	})
}
