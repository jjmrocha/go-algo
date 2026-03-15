package sort

// Quick sorts arr in place using the order defined by cmp.
func Quick[T any](arr []T, cmp Comparator[T]) {
	if len(arr) <= 1 {
		return
	}

	first := 0
	last := len(arr) - 1
	pivot := pivot(arr, cmp)

	i := 0

	for i <= last {
		switch cmp(arr[i], pivot) {
		case Before:
			Swap(arr, first, i)
			first++
			i++
		case After:
			Swap(arr, i, last)
			last--
		default:
			i++
		}
	}

	if first > 1 {
		Quick[T](arr[:first], cmp)
	}

	if last < len(arr)-1 {
		Quick[T](arr[last+1:], cmp)
	}
}

func pivot[T any](arr []T, cmp Comparator[T]) T {
	if len(arr) <= 3 {
		return arr[0]
	}

	lo, mid, hi := 0, len(arr)/2, len(arr)-1
	lv, mv, hv := arr[lo], arr[mid], arr[hi]

	if cmp(lv, mv) == Before {
		if cmp(mv, hv) == Before {
			return mv
		}

		if cmp(lv, hv) == Before {
			return hv
		}

		return lv
	}

	if cmp(lv, hv) == Before {
		return lv
	}

	if cmp(mv, hv) == Before {
		return hv
	}

	return mv
}
