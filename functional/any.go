package functional

import (
	"iter"
	"slices"
)

// Any reports whether fn returns true for at least one element of input.
func Any[T any](input []T, fn func(T) bool) bool {
	return slices.ContainsFunc(input, fn)
}

// AnySeq reports whether fn returns true for at least one element yielded by seq.
// It short-circuits on the first matching element.
func AnySeq[T any](seq iter.Seq[T], fn func(T) bool) bool {
	for v := range seq {
		if fn(v) {
			return true
		}
	}

	return false
}
