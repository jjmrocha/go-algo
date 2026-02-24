// Package stack provides a generic LIFO (Last-In, First-Out) stack implementation
// backed by a singly-linked list. It supports any element type via Go generics.
package stack

// node is an internal linked-list element that holds a value and a pointer to
// the previous node in the chain.
type node[T any] struct {
	previous *node[T]
	data     T
}

// Stack is a generic LIFO data structure. The zero value is ready to use.
type Stack[T any] struct {
	last *node[T]
	size int64
}

// New returns an empty Stack ready for use.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds data to the top of the stack.
func (s *Stack[T]) Push(data T) {
	s.last = &node[T]{
		previous: s.last,
		data:     data,
	}
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

	value := s.last.data
	s.last = s.last.previous
	s.size--

	return value, true
}

// Peek returns the top element without removing it.
// The second return value is false when the stack is empty, in which case
// the zero value of T is returned.
func (s *Stack[T]) Peek() (T, bool) {
	if s.size == 0 {
		var zero T
		return zero, false
	}

	return s.last.data, true
}

// Len returns the number of elements currently in the stack.
func (s *Stack[T]) Len() int64 {
	return s.size
}

// Empty reports whether the stack contains no elements.
func (s *Stack[T]) Empty() bool {
	return s.size == 0
}
