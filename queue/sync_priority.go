package queue

import (
	"iter"
	"sync"

	"github.com/jjmrocha/go-algo/sort"
)

// SyncPQueue is a thread-safe wrapper around [PQueue] using a read/write mutex.
// Mutating operations (Enqueue, Dequeue) hold an exclusive lock; read-only
// operations (Peek, Len, Empty) hold a shared lock.
type SyncPQueue[T any] struct {
	mu sync.RWMutex
	q  *PQueue[T]
}

// NewSyncPriorityQueue returns an empty, thread-safe SyncPQueue ready for use.
func NewSyncPriorityQueue[T any](cmp sort.Comparator[T]) *SyncPQueue[T] {
	return &SyncPQueue[T]{
		q: NewPriorityQueue[T](cmp),
	}
}

// NewSyncPriorityQueueWithCap returns an empty SyncPQueue with the given initial
// capacity. Returns an error if initialCap is not positive.
func NewSyncPriorityQueueWithCap[T any](initialCap int, cmp sort.Comparator[T]) (*SyncPQueue[T], error) {
	q, err := NewPriorityQueueWithCap[T](initialCap, cmp)
	if err != nil {
		return nil, err
	}
	return &SyncPQueue[T]{q: q}, nil
}

// Enqueue adds data to the queue. It is safe for concurrent use.
func (q *SyncPQueue[T]) Enqueue(data T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.q.Enqueue(data)
}

// Dequeue removes and returns the minimum element.
// The second return value is false when the queue is empty. It is safe for
// concurrent use.
func (q *SyncPQueue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.q.Dequeue()
}

// Peek returns the minimum element without removing it.
// The second return value is false when the queue is empty. It is safe for
// concurrent use.
func (q *SyncPQueue[T]) Peek() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Peek()
}

// Len returns the number of elements currently in the queue. It is safe for
// concurrent use.
func (q *SyncPQueue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Len()
}

// Empty reports whether the queue contains no elements. It is safe for
// concurrent use.
func (q *SyncPQueue[T]) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.q.Empty()
}

// Values returns an iterator that drains the queue in priority order.
// Each dequeued element is yielded once; the queue is empty after the
// iterator is fully consumed. It is safe for concurrent use.
func (q *SyncPQueue[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, ok := q.Dequeue(); ok; v, ok = q.Dequeue() {
			if !yield(v) {
				return
			}
		}
	}
}
