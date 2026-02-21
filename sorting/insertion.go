package sorting

import "golang.org/x/exp/constraints"

// InsertionSort sorts a slice using the insertion sort algorithm
func InsertionSort[T constraints.Ordered](arr []T) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}
