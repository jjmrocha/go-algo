package functional

import "iter"

// Filter returns a new slice containing only the elements of input for which
// fn returns true, preserving the original order.
func Filter[T any](input []T, fn func(T) bool) []T {
	result := make([]T, 0, len(input))

	for _, v := range input {
		if fn(v) {
			result = append(result, v)
		}
	}

	return result
}

// FilterSeq returns an iterator that yields only the elements of seq for which
// fn returns true, preserving order.
func FilterSeq[T any](seq iter.Seq[T], fn func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if fn(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
