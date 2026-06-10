package fn

import "iter"

// Distinct returns a new slice containing the elements of input with
// duplicates removed, preserving the first occurrence of each element.
func Distinct[T comparable](input []T) []T {
	seen := make(map[T]struct{}, len(input))
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// DistinctSeq returns an iterator that yields the elements of seq with
// duplicates removed, preserving the first occurrence of each element.
// Each new iteration starts with a fresh seen set.
func DistinctSeq[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[T]struct{}, defaultSize)

		for v := range seq {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				if !yield(v) {
					return
				}
			}
		}
	}
}
