package queue

import "errors"

// ErrCapacityTooSmall is returned by [NewBlockingQueue] and
// [NewPriorityQueueWithCap] (and its sync variant) when the requested capacity
// is not positive.
var ErrCapacityTooSmall = errors.New("capacity must be greater than zero")
