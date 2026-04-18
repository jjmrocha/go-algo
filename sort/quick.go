package sort

import "math/rand/v2"

// Quick sorts arr in place using the order defined by cmp.
func Quick[T any](arr []T, cmp Comparator[T]) {
	if len(arr) <= 1 {
		return
	}

	quick(arr, cmp)
}

func quick[T any](arr []T, cmp Comparator[T]) {
	if len(arr) <= 1 {
		return
	}

	first, last := partition(arr, cmp)
	quick[T](arr[:first], cmp)
	quick[T](arr[last+1:], cmp)
}

func partition[T any](arr []T, cmp Comparator[T]) (int, int) {
	pivot := rand.IntN(len(arr)) //nolint:gosec // G404: pivot choice needs speed, not cryptographic randomness
	Swap(arr, 0, pivot)

	partition := arr[0]

	i, lt, gt := 1, 0, len(arr)-1

	for i <= gt {
		switch cmp(arr[i], partition) {
		case Before:
			Swap(arr, lt, i)
			lt++
			i++
		case Equal:
			i++
		default:
			Swap(arr, i, gt)
			gt--
		}
	}

	return lt, gt
}
