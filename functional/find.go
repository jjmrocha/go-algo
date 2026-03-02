package functional

// Find returns the first element of input for which fn returns true, along
// with true. If no element matches, it returns the zero value of T and false.
func Find[T any](input []T, fn func(T) bool) (T, bool) {
	for _, v := range input {
		if fn(v) {
			return v, true
		}
	}

	var zero T
	return zero, false
}
