package queue

import "iter"

const defaultCap = 16

type PriorityQueue[T any] struct {
	items []T
	size  int
	cmp   func(a, b T) int
}

func NewPriorityQueue[T any](cmp func(a, b T) int) *PriorityQueue[T] {
	q, _ := NewPriorityQueueWithCap[T](defaultCap, cmp)
	return q
}

func NewPriorityQueueWithCap[T any](initialCap int, cmp func(a, b T) int) (*PriorityQueue[T], error) {
	if initialCap <= 0 {
		return nil, ErrCapacityTooSmall
	}

	return &PriorityQueue[T]{
		items: make([]T, initialCap),
		size:  0,
		cmp:   cmp,
	}, nil
}

func (q *PriorityQueue[T]) Enqueue(data T) {
	q.resizeIfNeeded()

	q.items[q.size] = data
	q.size++

	q.swim()
}

func (q *PriorityQueue[T]) Peek() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	return q.items[0], true
}

func (q *PriorityQueue[T]) Dequeue() (T, bool) {
	var zero T

	if q.size == 0 {
		return zero, false
	}

	data := q.items[0]
	q.size--
	q.swap(0, q.size)
	q.items[q.size] = zero

	q.sink()
	q.resizeIfNeeded()

	return data, true
}

func (q *PriorityQueue[T]) Len() int {
	return q.size
}

func (q *PriorityQueue[T]) Empty() bool {
	return q.size == 0
}

func (q *PriorityQueue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := q.Dequeue(); ok; v, ok = q.Dequeue() {
			if !yield(v) {
				return
			}
		}
	}
}

func (q *PriorityQueue[T]) resizeIfNeeded() {
	s := len(q.items)

	if q.size == s {
		q.resize(s * 2)
		return
	}

	if q.size < s/4 && s > defaultCap {
		q.resize(s / 2)
		return
	}
}

func (q *PriorityQueue[T]) resize(newCap int) {
	newItems := make([]T, newCap)
	copy(newItems, q.items[:q.size])
	q.items = newItems
}

func (q *PriorityQueue[T]) swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

func (q *PriorityQueue[T]) swim() {
	i := q.size - 1
	p := parent(i)

	for i > 0 && q.cmp(q.items[i], q.items[p]) < 0 {
		q.swap(i, p)

		i = p
		p = parent(i)
	}
}

func (q *PriorityQueue[T]) sink() {
	i := 0

	for {
		l := left(i)
		r := right(i)
		smallest := i

		if l < q.size && q.cmp(q.items[l], q.items[smallest]) < 0 {
			smallest = l
		}

		if r < q.size && q.cmp(q.items[r], q.items[smallest]) < 0 {
			smallest = r
		}

		if smallest == i {
			break
		}

		q.swap(i, smallest)
		i = smallest
	}
}

func parent(k int) int {
	return (k - 1) / 2
}

func left(k int) int {
	return (2 * k) + 1
}

func right(k int) int {
	return (2 * k) + 2
}
