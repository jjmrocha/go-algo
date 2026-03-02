package functional

// Fold reduces input to a single value by applying fn to an accumulator and
// each element from left to right, starting with initial as the accumulator.
func Fold[T any, U any](input []T, initial U, fn func(U, T) U) U {
	result := initial

	for _, v := range input {
		result = fn(result, v)
	}

	return result
}
