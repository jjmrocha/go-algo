package sorting

import "golang.org/x/exp/constraints"

// QuickSort sorts a slice using the quicksort algorithm
func QuickSort[T constraints.Ordered](arr []T) {
	if len(arr) <= 1 {
		return
	}
	quickSort(arr, 0, len(arr)-1)
}

func quickSort[T constraints.Ordered](arr []T, lo, hi int) {
	if hi <= lo {
		return
	}
	j := partition(arr, lo, hi)
	quickSort(arr, lo, j-1)
	quickSort(arr, j+1, hi)
}

func partition[T constraints.Ordered](arr []T, lo, hi int) int {
	i := lo
	j := hi + 1
	pivot := arr[lo]

	for {
		// Find item on left to swap
		i++
		for i <= hi && arr[i] < pivot {
			i++
		}

		// Find item on right to swap
		j--
		for j >= lo && arr[j] > pivot {
			j--
		}

		// Check if pointers cross
		if i >= j {
			break
		}

		// Swap
		arr[i], arr[j] = arr[j], arr[i]
	}

	// Swap with partitioning item
	arr[lo], arr[j] = arr[j], arr[lo]
	return j
}
