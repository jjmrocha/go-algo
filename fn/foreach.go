package fn

import "iter"

// ForEach calls fn for each element of input in order.
// It is intended for side-effecting operations; use [Map] when a result is needed.
func ForEach[T any](input []T, fn func(T)) {
	for _, v := range input {
		fn(v)
	}
}

// ForEachSeq calls fn for each element yielded by seq in order.
// It is intended for side-effecting operations; use [MapSeq] when a result is needed.
func ForEachSeq[T any](seq iter.Seq[T], fn func(T)) {
	for v := range seq {
		fn(v)
	}
}
