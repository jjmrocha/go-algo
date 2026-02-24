package queue

import "sync"

// SyncQueue is a thread-safe wrapper around [Queue] that uses a read/write
// mutex to allow concurrent reads while serialising writes.
type SyncQueue[T any] struct {
	mu sync.RWMutex
	q  *Queue[T]
}

// NewSyncQueue returns an empty, thread-safe SyncQueue ready for use.
func NewSyncQueue[T any]() *SyncQueue[T] {
	return &SyncQueue[T]{
		q: New[T](),
	}
}

// Enqueue adds v to the back of the queue. It is safe for concurrent use.
func (q *SyncQueue[T]) Enqueue(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Enqueue(v)
}

// Dequeue removes and returns the front element of the queue.
// The second return value is false when the queue is empty, in which case
// the zero value of T is returned. It is safe for concurrent use.
func (q *SyncQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Dequeue()
}

// Len returns the number of elements currently in the queue.
// It is safe for concurrent use.
func (q *SyncQueue[T]) Len() int64 {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Len()
}

// Empty reports whether the queue contains no elements.
// It is safe for concurrent use.
func (q *SyncQueue[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Empty()
}
