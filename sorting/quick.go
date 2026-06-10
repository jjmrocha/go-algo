package sorting

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
	// middle-index pivot; swapped to position 0 before 3-way partitioning
	pivot := len(arr) / 2
	Swap(arr, 0, pivot)

	pivotVal := arr[0]

	i, lt, gt := 1, 0, len(arr)-1

	for i <= gt {
		switch c := cmp(arr[i], pivotVal); {
		case c < 0:
			Swap(arr, lt, i)
			lt++
			i++
		case c == 0:
			i++
		default:
			Swap(arr, i, gt)
			gt--
		}
	}

	return lt, gt
}
