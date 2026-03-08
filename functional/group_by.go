package functional

import "iter"

// GroupBy groups the elements of input by the key produced by fn.
// Elements that map to the same key are collected into a slice under that key,
// preserving their original order within each group.
func GroupBy[T any, K comparable](input []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T, defaultSize)

	for _, v := range input {
		key := fn(v)
		result[key] = append(result[key], v)
	}

	return result
}

// GroupBySeq groups the elements of seq by the key produced by fn and returns
// an iterator of (key, group) pairs. Each group preserves the order of elements
// as they appear in seq. The group map is rebuilt on each new iteration.
func GroupBySeq[T any, K comparable](seq iter.Seq[T], fn func(T) K) iter.Seq2[K, []T] {
	return func(yield func(K, []T) bool) {
		groups := make(map[K][]T, defaultSize)

		for v := range seq {
			key := fn(v)
			groups[key] = append(groups[key], v)
		}

		for k, group := range groups {
			if !yield(k, group) {
				return
			}
		}
	}
}
