package functional

import "slices"

// Any reports whether fn returns true for at least one element of input.
func Any[T any](input []T, fn func(T) bool) bool {
	return slices.ContainsFunc(input, fn)
}
