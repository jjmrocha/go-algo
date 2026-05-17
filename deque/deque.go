// Package deque provides a generic double-ended queue (deque) backed by a
// dynamically resized ring buffer. It supports efficient O(1) amortised push
// and pop at both ends and any element type via Go generics.
package deque

import (
	"errors"
	"iter"
)

const defaultCap = 16

// ErrorCapacityLessThanZero is returned when an attempt is made to create a Deque with a negative initial capacity.
var ErrCapacityTooSmall = errors.New("capacity must be greater than zero")

// Deque is a generic double-ended queue backed by a ring buffer.
// Use New to create a Deque; the zero value is not ready for use.
// A Deque must not be copied after first use.
type Deque[T any] struct {
	items []T
	first int
	last  int
	size  int
	cap   int
}

// New returns an empty Deque ready for use.
func New[T any]() *Deque[T] {
	d, _ := NewWithCap[T](defaultCap)
	return d
}

// NewWithCap returns an empty Deque with the specified initial capacity.
// It returns ErrCapacityTooSmall if capacity is not positive.
func NewWithCap[T any](initialCap int) (*Deque[T], error) {
	if initialCap <= 0 {
		return nil, ErrCapacityTooSmall
	}

	return &Deque[T]{
		items: make([]T, initialCap),
		first: 0,
		last:  0,
		size:  0,
		cap:   initialCap,
	}, nil
}

// PushFront adds data to the front of the deque.
func (d *Deque[T]) PushFront(data T) {
	d.resizeIfNeeded()

	d.first = d.prev(d.first)
	d.items[d.first] = data
	d.size++
}

// PushBack adds data to the back of the deque.
func (d *Deque[T]) PushBack(data T) {
	d.resizeIfNeeded()

	d.items[d.last] = data
	d.last = d.next(d.last)
	d.size++
}

// PopFront removes and returns the front element of the deque.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned.
func (d *Deque[T]) PopFront() (T, bool) {
	var zero T

	if d.size == 0 {
		return zero, false
	}

	data := d.items[d.first]
	d.items[d.first] = zero

	d.first = d.next(d.first)
	d.size--

	d.resizeIfNeeded()

	return data, true
}

// PopBack removes and returns the back element of the deque.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned.
func (d *Deque[T]) PopBack() (T, bool) {
	var zero T

	if d.size == 0 {
		return zero, false
	}

	d.last = d.prev(d.last)

	data := d.items[d.last]
	d.items[d.last] = zero

	d.size--

	d.resizeIfNeeded()

	return data, true
}

// PeekFront returns the front element without removing it.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned.
func (d *Deque[T]) PeekFront() (T, bool) {
	var zero T

	if d.size == 0 {
		return zero, false
	}

	return d.items[d.first], true
}

// PeekBack returns the back element without removing it.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned.
func (d *Deque[T]) PeekBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}

	i := d.prev(d.last)
	return d.items[i], true
}

// Len returns the number of elements currently in the deque.
func (d *Deque[T]) Len() int {
	return d.size
}

// Empty reports whether the deque contains no elements.
func (d *Deque[T]) Empty() bool {
	return d.size == 0
}

// Values returns an iterator over all elements from front to back.
// The deque is not modified; elements are not removed during iteration.
func (d *Deque[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := d.first

		for range d.size {
			if !yield(d.items[current]) {
				return
			}

			current = d.next(current)
		}
	}
}

// ToSlice returns a new slice containing all elements from front to back.
func (d *Deque[T]) ToSlice() []T {
	slice := make([]T, d.size)

	current := d.first

	for i := range d.size {
		slice[i] = d.items[current]
		current = d.next(current)
	}

	return slice
}

func (d *Deque[T]) resizeIfNeeded() {
	if d.size == d.cap {
		d.resize(d.cap * 2)
		return
	}

	if d.size < d.cap/4 && d.cap > defaultCap {
		d.resize(d.cap / 2)
		return
	}
}

func (d *Deque[T]) resize(newCap int) {
	newItems := make([]T, newCap)

	current := d.first

	for i := range d.size {
		newItems[i] = d.items[current]
		current = d.next(current)
	}

	d.items = newItems
	d.first = 0
	d.last = d.size
	d.cap = newCap
}

func (d *Deque[T]) next(i int) int {
	n := i + 1

	if n >= d.cap {
		return 0
	}

	return n
}

func (d *Deque[T]) prev(i int) int {
	p := i - 1

	if p < 0 {
		return d.cap - 1
	}

	return p
}
