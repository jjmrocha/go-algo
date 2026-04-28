package queue

// BlockingQueue is a bounded, channel-backed FIFO queue. Enqueue blocks
// when the queue is full; Dequeue blocks when it is empty. BlockingQueue
// is safe for concurrent use by multiple goroutines.
type BlockingQueue[T any] struct {
	ch chan T
}

// NewBlockingQueue creates a BlockingQueue with the given positive capacity.
// It returns ErrCapacityTooSmall if capacity is not positive.
func NewBlockingQueue[T any](capacity int) (*BlockingQueue[T], error) {
	if capacity <= 0 {
		return nil, ErrCapacityTooSmall
	}

	return &BlockingQueue[T]{ch: make(chan T, capacity)}, nil
}

// Enqueue adds v to the back of the queue. It blocks if the queue is full.
func (q *BlockingQueue[T]) Enqueue(v T) {
	q.ch <- v
}

// Dequeue removes and returns the front element of the queue.
// It blocks until an element is available if the queue is empty.
func (q *BlockingQueue[T]) Dequeue() T {
	return <-q.ch
}

// Len returns the number of elements currently in the queue.
func (q *BlockingQueue[T]) Len() int {
	return len(q.ch)
}

// Empty reports whether the queue contains no elements.
func (q *BlockingQueue[T]) Empty() bool {
	return len(q.ch) == 0
}

// Cap returns the maximum capacity of the queue.
func (q *BlockingQueue[T]) Cap() int {
	return cap(q.ch)
}

// Full reports whether the queue has reached its maximum capacity.
func (q *BlockingQueue[T]) Full() bool {
	return len(q.ch) == cap(q.ch)
}
