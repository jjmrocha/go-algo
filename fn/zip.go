package fn

// Zip combines elements from a and b pairwise using fn, returning a new slice
// of the results. The length of the result equals the length of the shorter
// input; excess elements from the longer slice are ignored.
func Zip[T any, U any, V any](a []T, b []U, fn func(T, U) V) []V {
	n := min(len(a), len(b))
	result := make([]V, n)

	for i := range n {
		result[i] = fn(a[i], b[i])
	}

	return result
}
