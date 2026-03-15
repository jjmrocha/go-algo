package sort

// Order represents the relative ordering between two values as returned by a Comparator.
type Order int

const (
	// Before indicates the first value is less than the second.
	Before Order = iota - 1
	// Equal indicates the two values are identical in ordering.
	Equal
	// After indicates the first value is greater than the second.
	After
)

// Comparator is a function that compares two values of type T and returns an Order indicating their relative order.
// It returns Before if a is less than b, Equal if a is equal to b, and After if a is greater than b.
type Comparator[T any] func(a, b T) Order

// Swap exchanges the elements at indices i and j in arr.
func Swap[T any](arr []T, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
