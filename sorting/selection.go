package sorting

import "golang.org/x/exp/constraints"

// SelectionSort sorts a slice using the selection sort algorithm
func SelectionSort[T constraints.Ordered](arr []T) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[i], arr[minIdx] = arr[minIdx], arr[i]
	}
}
