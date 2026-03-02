package functional

// ForEach calls fn for each element of input in order.
// It is intended for side-effecting operations; use [Map] when a result is needed.
func ForEach[T any](input []T, fn func(T)) {
	for _, v := range input {
		fn(v)
	}
}
