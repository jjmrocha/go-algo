package functional

// All reports whether fn returns true for every element of input.
// It short-circuits on the first non-matching element and returns true for an
// empty slice (vacuous truth).
func All[T any](input []T, fn func(T) bool) bool {
	for _, v := range input {
		if !fn(v) {
			return false
		}
	}

	return true
}
