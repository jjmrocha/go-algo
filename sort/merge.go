package sort

// Merge implements the merge sort algorithm.
func Merge[T any](arr []T, cmp Comparator[T]) {
	if len(arr) <= 1 {
		return
	}

	aux := make([]T, len(arr))
	sort[T](arr, aux, 0, len(arr), cmp)
}

func sort[T any](arr, aux []T, lo, hi int, cmp Comparator[T]) {
	if hi-lo <= 1 {
		return
	}

	mid := lo + (hi-lo)/2
	sort[T](arr, aux, lo, mid, cmp)
	sort[T](arr, aux, mid, hi, cmp)

	merge(arr, aux, lo, mid, hi, cmp)
}

func merge[T any](arr, aux []T, lo, mid, hi int, cmp Comparator[T]) {
	for i := lo; i < hi; i++ {
		aux[i] = arr[i]
	}

	i, j := lo, mid

	for k := lo; k < hi; k++ {
		if i >= mid {
			arr[k] = aux[j]
			j++
		} else if j >= hi {
			arr[k] = aux[i]
			i++
		} else if cmp(aux[i], aux[j]) == Before {
			arr[k] = aux[i]
			i++
		} else {
			arr[k] = aux[j]
			j++
		}
	}
}
