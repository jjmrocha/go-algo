package sorting

// Insertion sorts arr in place using the order defined by cmp.
func Insertion[T any](arr []T, cmp Comparator[T]) {
	insertionWithStep(arr, cmp, 1)
}

func insertionWithStep[T any](arr []T, cmp Comparator[T], step int) {
	for i := step; i < len(arr); i++ {
		for j := i; j >= step; j -= step {
			if cmp(arr[j-step], arr[j]) <= 0 {
				break
			}

			Swap(arr, j-step, j)
		}
	}
}
