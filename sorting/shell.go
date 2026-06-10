package sorting

var shellSteps = []int{701, 301, 132, 57, 23, 10, 4, 1}

// Shell sorts arr in place using the order defined by cmp, using Ciura's gap sequence.
func Shell[T any](arr []T, cmp Comparator[T]) {
	l := len(arr)

	for _, step := range shellSteps {
		if step >= l {
			continue
		}

		insertionWithStep(arr, cmp, step)
	}
}
