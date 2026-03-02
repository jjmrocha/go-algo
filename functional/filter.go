package functional

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
