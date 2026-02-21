package sorting

import "golang.org/x/exp/constraints"

// MergeSort sorts a slice using the merge sort algorithm
func MergeSort[T constraints.Ordered](arr []T) {
	if len(arr) <= 1 {
		return
	}
	aux := make([]T, len(arr))
	mergeSort(arr, aux, 0, len(arr)-1)
}

func mergeSort[T constraints.Ordered](arr, aux []T, lo, hi int) {
	if hi <= lo {
		return
	}
	mid := lo + (hi-lo)/2
	mergeSort(arr, aux, lo, mid)
	mergeSort(arr, aux, mid+1, hi)
	merge(arr, aux, lo, mid, hi)
}

func merge[T constraints.Ordered](arr, aux []T, lo, mid, hi int) {
	// Copy to aux
	for k := lo; k <= hi; k++ {
		aux[k] = arr[k]
	}

	// Merge back to arr
	i := lo
	j := mid + 1
	for k := lo; k <= hi; k++ {
		if i > mid {
			arr[k] = aux[j]
			j++
		} else if j > hi {
			arr[k] = aux[i]
			i++
		} else if aux[j] < aux[i] {
			arr[k] = aux[j]
			j++
		} else {
			arr[k] = aux[i]
			i++
		}
	}
}
