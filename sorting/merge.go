package sorting

// Merge sorts arr in place using the order defined by cmp. It is a stable sort
// running in O(n log n) time with a single auxiliary buffer allocation.
func Merge[T any](arr []T, cmp Comparator[T]) {
	if len(arr) <= 1 {
		return
	}

	aux := make([]T, len(arr))
	mergeSort[T](arr, aux, 0, len(arr), cmp)
}

func mergeSort[T any](arr, aux []T, lo, hi int, cmp Comparator[T]) {
	if hi-lo <= 1 {
		return
	}

	mid := lo + (hi-lo)/2
	mergeSort[T](arr, aux, lo, mid, cmp)
	mergeSort[T](arr, aux, mid, hi, cmp)

	merge(arr, aux, lo, mid, hi, cmp)
}

func merge[T any](arr, aux []T, lo, mid, hi int, cmp Comparator[T]) {
	for i := lo; i < hi; i++ {
		aux[i] = arr[i]
	}

	i, j := lo, mid

	for k := lo; k < hi; k++ {
		switch {
		case i >= mid:
			arr[k] = aux[j]
			j++
		case j >= hi:
			arr[k] = aux[i]
			i++
		case cmp(aux[i], aux[j]) < 0:
			arr[k] = aux[i]
			i++
		default:
			arr[k] = aux[j]
			j++
		}
	}
}
