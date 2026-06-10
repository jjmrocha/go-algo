// Package queue provides a generic FIFO (First-In, First-Out) queue
// implementation backed by a singly-linked list. It supports any element
// type via Go generics.
package queue

import (
	"iter"
	"sync"
)

type node[T any] struct {
	next *node[T]
	data T
}

// Queue is a generic FIFO data structure.
// Use [New] to create a Queue. A Queue must not be copied after first use.
type Queue[T any] struct {
	pool  sync.Pool
	first *node[T]
	last  *node[T]
	size  int
}

// New returns an empty Queue ready for use.
func New[T any]() *Queue[T] {
	q := Queue[T]{}
	q.pool.New = func() any {
		return &node[T]{}
	}

	return &q
}

// Enqueue adds data to the back of the queue.
func (q *Queue[T]) Enqueue(data T) {
	n := q.pool.Get().(*node[T])
	n.data = data
	n.next = nil

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

	n := q.first
	q.first = n.next
	q.size--

	if q.size == 0 {
		q.last = nil // clear stale tail pointer when queue becomes empty
	}

	val := n.data
	q.pool.Put(n)

	return val, true
}

// Len returns the number of elements currently in the queue.
func (q *Queue[T]) Len() int {
	return q.size
}

// Empty reports whether the queue contains no elements.
func (q *Queue[T]) Empty() bool {
	return q.size == 0
}

// Drain returns an iterator that drains the queue from front to back.
// Each dequeued element is yielded once; the queue is empty after the
// iterator is fully consumed.
func (q *Queue[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := q.Dequeue(); ok; v, ok = q.Dequeue() {
			if !yield(v) {
				return
			}
		}
	}
}
