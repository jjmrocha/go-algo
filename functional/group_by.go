package functional

// GroupBy groups the elements of input by the key produced by fn.
// Elements that map to the same key are collected into a slice under that key,
// preserving their original order within each group.
func GroupBy[T any, K comparable](input []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)

	for _, v := range input {
		key := fn(v)
		result[key] = append(result[key], v)
	}

	return result
}
