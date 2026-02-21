package sorting

import "golang.org/x/exp/constraints"

// QuickSelect finds the kth smallest element in the array
// k is 0-indexed (k=0 returns the smallest element)
func QuickSelect[T constraints.Ordered](arr []T, k int) T {
	if k < 0 || k >= len(arr) {
		panic("k out of bounds")
	}
	return quickSelect(arr, 0, len(arr)-1, k)
}

func quickSelect[T constraints.Ordered](arr []T, lo, hi, k int) T {
	if lo == hi {
		return arr[lo]
	}

	j := partition(arr, lo, hi)

	if j == k {
		return arr[k]
	} else if j > k {
		return quickSelect(arr, lo, j-1, k)
	} else {
		return quickSelect(arr, j+1, hi, k)
	}
}
