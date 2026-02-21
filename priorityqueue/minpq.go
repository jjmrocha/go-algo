package priorityqueue

import "golang.org/x/exp/constraints"

// MinPQ is a minimum priority queue implementation using a binary heap
type MinPQ[T constraints.Ordered] struct {
	pq []T
	n  int
}

// NewMinPQ creates a new minimum priority queue
func NewMinPQ[T constraints.Ordered]() *MinPQ[T] {
	return &MinPQ[T]{
		pq: make([]T, 1), // index 0 is not used
		n:  0,
	}
}

// Insert adds a new element to the priority queue
func (pq *MinPQ[T]) Insert(item T) {
	pq.n++
	if pq.n >= len(pq.pq) {
		pq.pq = append(pq.pq, item)
	} else {
		pq.pq[pq.n] = item
	}
	pq.swim(pq.n)
}

// DelMin removes and returns the minimum element
func (pq *MinPQ[T]) DelMin() (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}
	min := pq.pq[1]
	pq.pq[1], pq.pq[pq.n] = pq.pq[pq.n], pq.pq[1]
	pq.n--
	pq.sink(1)
	return min, true
}

// Min returns the minimum element without removing it
func (pq *MinPQ[T]) Min() (T, bool) {
	if pq.IsEmpty() {
		var zero T
		return zero, false
	}
	return pq.pq[1], true
}

// IsEmpty returns true if the priority queue is empty
func (pq *MinPQ[T]) IsEmpty() bool {
	return pq.n == 0
}

// Size returns the number of elements in the priority queue
func (pq *MinPQ[T]) Size() int {
	return pq.n
}

func (pq *MinPQ[T]) swim(k int) {
	for k > 1 && pq.pq[k/2] > pq.pq[k] {
		pq.pq[k/2], pq.pq[k] = pq.pq[k], pq.pq[k/2]
		k = k / 2
	}
}

func (pq *MinPQ[T]) sink(k int) {
	for 2*k <= pq.n {
		j := 2 * k
		if j < pq.n && pq.pq[j] > pq.pq[j+1] {
			j++
		}
		if pq.pq[k] <= pq.pq[j] {
			break
		}
		pq.pq[k], pq.pq[j] = pq.pq[j], pq.pq[k]
		k = j
	}
}
