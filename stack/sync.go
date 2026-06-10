package stack

import (
	"iter"
	"sync"
)

// SyncStack is a thread-safe wrapper around [Stack] that uses a read/write
// mutex to allow concurrent reads while serialising writes.
type SyncStack[T any] struct {
	mu sync.RWMutex
	s  *Stack[T]
}

// NewSyncStack returns an empty, thread-safe SyncStack ready for use.
func NewSyncStack[T any]() *SyncStack[T] {
	return &SyncStack[T]{
		s: New[T](),
	}
}

// Push adds value to the top of the stack. It is safe for concurrent use.
func (s *SyncStack[T]) Push(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s.Push(value)
}

// Pop removes and returns the top element of the stack.
// The second return value is false when the stack is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (s *SyncStack[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.s.Pop()
}

// Peek returns the top element without removing it.
// The second return value is false when the stack is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (s *SyncStack[T]) Peek() (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Peek()
}

// Len returns the number of elements currently in the stack.
// It is safe for concurrent use.
func (s *SyncStack[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Len()
}

// Empty reports whether the stack contains no elements.
// It is safe for concurrent use.
func (s *SyncStack[T]) Empty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.s.Empty()
}

// Drain returns an iterator that drains the stack from top to bottom.
// Each popped element is yielded once; the stack is empty after the
// iterator is fully consumed. It is safe for concurrent use.
func (s *SyncStack[T]) Drain() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := s.Pop(); ok; v, ok = s.Pop() {
			if !yield(v) {
				return
			}
		}
	}
}
