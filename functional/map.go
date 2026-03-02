// Package functional provides generic higher-order functions for slices,
// following functional programming conventions.
//
// All functions are pure: they do not modify their input and always return
// new slices or values. nil inputs are treated as empty slices.
package functional

// Map applies fn to each element of input and returns a new slice containing
// the results in the same order.
func Map[T any, U any](input []T, fn func(T) U) []U {
	result := make([]U, len(input))

	for i, v := range input {
		result[i] = fn(v)
	}

	return result
}
