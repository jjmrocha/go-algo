package fn

// Partition splits input into two slices: the first contains all elements for
// which fn returns true, the second contains all elements for which fn returns
// false. The relative order of elements is preserved in each slice.
func Partition[T any](input []T, fn func(T) bool) ([]T, []T) {
	var matching, nonMatching []T

	for _, v := range input {
		if fn(v) {
			matching = append(matching, v)
		} else {
			nonMatching = append(nonMatching, v)
		}
	}

	return matching, nonMatching
}
