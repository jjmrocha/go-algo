package functional

import (
	"maps"
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
		// when
		result := GroupBy([]int(nil), parity)
		// then
		if len(result) != 0 {
			t.Errorf("GroupBy nil: got %v, want empty map", result)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		// when
		result := GroupBy([]int{}, parity)
		// then
		if len(result) != 0 {
			t.Errorf("GroupBy empty: got %v, want empty map", result)
		}
	})

	t.Run("all same key", func(t *testing.T) {
		// given
		input := []int{2, 4, 6}
		// when
		result := GroupBy(input, parity)
		// then
		if len(result) != 1 {
			t.Fatalf("expected 1 group, got %d", len(result))
		}
		if !slices.Equal(result["even"], []int{2, 4, 6}) {
			t.Errorf("GroupBy all-same: got %v, want [2 4 6]", result["even"])
		}
	})

	t.Run("all different keys", func(t *testing.T) {
		// given
		input := []int{1, 2, 3}
		// when
		result := GroupBy(input, func(v int) int { return v })
		// then
		if len(result) != 3 {
			t.Fatalf("expected 3 groups, got %d", len(result))
		}
		for _, v := range []int{1, 2, 3} {
			if !slices.Equal(result[v], []int{v}) {
				t.Errorf("GroupBy key %d: got %v, want [%d]", v, result[v], v)
			}
		}
	})

	t.Run("mixed — order within group preserved", func(t *testing.T) {
		// given
		input := []int{1, 3, 2, 5, 4}
		// when
		result := GroupBy(input, parity)
		// then
		if !slices.Equal(result["odd"], []int{1, 3, 5}) {
			t.Errorf("GroupBy odd order: got %v, want [1 3 5]", result["odd"])
		}
		if !slices.Equal(result["even"], []int{2, 4}) {
			t.Errorf("GroupBy even order: got %v, want [2 4]", result["even"])
		}
	})
}

func TestGroupBySeq(t *testing.T) {
	parity := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}

	t.Run("empty seq", func(t *testing.T) {
		// when
		result := maps.Collect(GroupBySeq(slices.Values([]int{}), parity))
		// then
		if len(result) != 0 {
			t.Errorf("GroupBySeq empty = %v, want empty map", result)
		}
	})

	t.Run("all same key", func(t *testing.T) {
		// given
		input := []int{2, 4, 6}
		// when
		result := maps.Collect(GroupBySeq(slices.Values(input), parity))
		// then
		if !slices.Equal(result["even"], []int{2, 4, 6}) {
			t.Errorf("GroupBySeq all-same even = %v, want [2 4 6]", result["even"])
		}
	})

	t.Run("mixed — order within group preserved", func(t *testing.T) {
		// given
		input := []int{1, 2, 3, 4, 5}
		// when
		result := maps.Collect(GroupBySeq(slices.Values(input), parity))
		// then
		if !slices.Equal(result["odd"], []int{1, 3, 5}) {
			t.Errorf("GroupBySeq odd group = %v, want [1 3 5]", result["odd"])
		}
		if !slices.Equal(result["even"], []int{2, 4}) {
			t.Errorf("GroupBySeq even group = %v, want [2 4]", result["even"])
		}
	})
}

// TestGroupBySeq_MultipleIterations verifies that each iteration rebuilds the
// group map — state must not accumulate across iterations.
func TestGroupBySeq_MultipleIterations(t *testing.T) {
	// given
	parity := func(v int) string {
		if v%2 == 0 {
			return "even"
		}
		return "odd"
	}
	seq := GroupBySeq(slices.Values([]int{1, 2, 3, 4}), parity)
	// when
	first := maps.Collect(seq)
	second := maps.Collect(seq)
	// then
	if !slices.Equal(first["odd"], second["odd"]) {
		t.Errorf("odd group: first=%v, second=%v — state leaked across iterations", first["odd"], second["odd"])
	}
	if !slices.Equal(first["even"], second["even"]) {
		t.Errorf("even group: first=%v, second=%v — state leaked across iterations", first["even"], second["even"])
	}
}
