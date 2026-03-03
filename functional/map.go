// Package functional provides generic higher-order functions for slices,
// following functional programming conventions.
//
// All functions are pure: they do not modify their input and always return
// new slices or values. nil inputs are treated as empty slices.
package functional

import "iter"

// Map applies fn to each element of input and returns a new slice containing
// the results in the same order.
func Map[T any, U any](input []T, fn func(T) U) []U {
	result := make([]U, len(input))

	for i, v := range input {
		result[i] = fn(v)
	}

	return result
}

// MapSeq returns an iterator that applies fn to each element of seq and yields
// the results in the same order. Transformation is lazy: fn is called only as
// the caller iterates the returned sequence.
func MapSeq[T any, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}
