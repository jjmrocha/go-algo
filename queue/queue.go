// Package queue provides a generic FIFO (First-In, First-Out) queue
// implementation backed by a singly-linked list. It supports any element
// type via Go generics.
package queue

type node[T any] struct {
	next *node[T]
	data T
}

// Queue is a generic FIFO data structure. The zero value is ready to use.
type Queue[T any] struct {
	first *node[T]
	last  *node[T]
	size  int64
}

// New returns an empty Queue ready for use.
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

// Enqueue adds data to the back of the queue.
func (q *Queue[T]) Enqueue(data T) {
	n := &node[T]{
		data: data,
	}

	if q.size == 0 {
		q.first = n
		q.last = n
	} else {
		q.last.next = n
		q.last = n
	}

	q.size++
}

// Dequeue removes and returns the front element of the queue.
// The second return value is false when the queue is empty, in which case
// the zero value of T is returned.
func (q *Queue[T]) Dequeue() (T, bool) {
	if q.size == 0 {
		var zero T
		return zero, false
	}

	v := q.first.data
	q.first = q.first.next
	q.size--

	return v, true
}

// Len returns the number of elements currently in the queue.
func (q *Queue[T]) Len() int64 {
	return q.size
}

// Empty reports whether the queue contains no elements.
func (q *Queue[T]) Empty() bool {
	return q.size == 0
}
