// Package stack provides a generic LIFO (Last-In, First-Out) stack implementation
// backed by a singly-linked list. It supports any element type via Go generics.
package stack

import (
	"iter"
	"sync"
)

// node is an internal linked-list element that holds a value and a pointer to
// the next node in the chain.
type node[T any] struct {
	next *node[T]
	data T
}

// Stack is a generic LIFO data structure.
// Use [New] to create a Stack. A Stack must not be copied after first use.
type Stack[T any] struct {
	pool  sync.Pool
	first *node[T]
	size  int
}

// New returns an empty Stack ready for use.
func New[T any]() *Stack[T] {
	s := Stack[T]{}
	s.pool.New = func() any {
		return &node[T]{}
	}

	return &s
}

// Push adds data to the top of the stack.
func (s *Stack[T]) Push(data T) {
	n := s.pool.Get().(*node[T])
	n.data = data
	n.next = s.first
	s.first = n
	s.size++
}

// Pop removes and returns the top element of the stack.
// The second return value is false when the stack is empty, in which case
// the zero value of T is returned.
func (s *Stack[T]) Pop() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}

	n := s.first
	s.first = n.next
	s.size--

	val := n.data
	s.pool.Put(n)

	return val, true
}

// Peek returns the top element without removing it.
// The second return value is false when the stack is empty, in which case
// the zero value of T is returned.
func (s *Stack[T]) Peek() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}

	return s.first.data, true
}

// Len returns the number of elements currently in the stack.
func (s *Stack[T]) Len() int {
	return s.size
}

// Empty reports whether the stack contains no elements.
func (s *Stack[T]) Empty() bool {
	return s.size == 0
}

// Drain returns an iterator that drains the stack from top to bottom.
// Each popped element is yielded once; the stack is empty after the
// iterator is fully consumed.
func (s *Stack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := s.Pop(); ok; v, ok = s.Pop() {
			if !yield(v) {
				return
			}
		}
	}
}
