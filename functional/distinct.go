package functional

// Distinct returns a new slice containing the elements of input with
// duplicates removed, preserving the first occurrence of each element.
func Distinct[T comparable](input []T) []T {
	seen := make(map[T]struct{}, len(input))
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}
