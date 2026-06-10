// Package sorting provides generic, in-place sorting and shuffling algorithms
// parameterised by a [Comparator]. The comparator follows the standard library
// convention used by [cmp.Compare] and [slices.SortFunc]: it returns a negative
// number when a is less than b, zero when they are equal, and a positive number
// when a is greater than b.
package sorting

// Comparator is a function that compares two values of type T.
// It returns a negative number if a is less than b, zero if a equals b, and a
// positive number if a is greater than b — the same contract as [cmp.Compare].
type Comparator[T any] func(a, b T) int

// Swap exchanges the elements at indices i and j in arr.
func Swap[T any](arr []T, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
