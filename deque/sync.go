package deque

import (
	"iter"
	"slices"
	"sync"
)

// SyncDeque is a thread-safe wrapper around [Deque] that uses a read/write
// mutex to allow concurrent reads while serialising writes.
type SyncDeque[T any] struct {
	mu sync.RWMutex
	d  *Deque[T]
}

// NewSyncDeque returns an empty, thread-safe SyncDeque ready for use.
func NewSyncDeque[T any]() *SyncDeque[T] {
	return &SyncDeque[T]{
		d: New[T](),
	}
}

// PushFront adds v to the front of the deque. It is safe for concurrent use.
func (d *SyncDeque[T]) PushFront(v T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.d.PushFront(v)
}

// PushBack adds v to the back of the deque. It is safe for concurrent use.
func (d *SyncDeque[T]) PushBack(v T) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.d.PushBack(v)
}

// PopFront removes and returns the front element of the deque.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (d *SyncDeque[T]) PopFront() (T, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.d.PopFront()
}

// PopBack removes and returns the back element of the deque.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (d *SyncDeque[T]) PopBack() (T, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.d.PopBack()
}

// PeekFront returns the front element without removing it.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (d *SyncDeque[T]) PeekFront() (T, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.d.PeekFront()
}

// PeekBack returns the back element without removing it.
// The second return value is false when the deque is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (d *SyncDeque[T]) PeekBack() (T, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.d.PeekBack()
}

// Len returns the number of elements currently in the deque.
// It is safe for concurrent use.
func (d *SyncDeque[T]) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.d.Len()
}

// Empty reports whether the deque contains no elements.
// It is safe for concurrent use.
func (d *SyncDeque[T]) Empty() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.d.Empty()
}

// Values returns an iterator over all elements from front to back.
// It takes a snapshot of the deque and is safe for concurrent use.
func (d *SyncDeque[T]) Values() iter.Seq[T] {
	return slices.Values(d.ToSlice())
}

// ToSlice returns a new slice containing all elements from front to back.
// It is safe for concurrent use.
func (d *SyncDeque[T]) ToSlice() []T {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.d.ToSlice()
}
