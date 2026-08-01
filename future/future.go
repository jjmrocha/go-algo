// Package future provides a generic Future/Promise implementation for Go,
// enabling asynchronous computation with context-aware cancellation.
//
// A Future represents a value that will be available at some point in the
// future, produced by an asynchronous operation. Multiple callers may Await
// the same Future concurrently — all are unblocked simultaneously when the
// result is ready, and all receive the same value and error.
//
// Example usage:
//
//	f := future.AsyncWithContext(ctx, func(ctx context.Context) (int, error) {
//	    return expensiveComputation(ctx)
//	})
//
//	result, err := f.AwaitWithContext(ctx)
package future

import (
	"context"
	"sync"
	"time"
)

// Future represents the result of an asynchronous computation of type T.
// Multiple concurrent calls to [Future.AwaitWithContext] are safe: all callers block
// until the result is ready and then all receive the same value and error.
// The result is cached, so calls after the computation completes return immediately.
type Future[T any] struct {
	once sync.Once
	done chan struct{}
	val  T
	err  error
}

// New creates an unresolved Future. The caller is responsible for completing
// it with [Future.Resolve]; until then all Await variants block. Futures must
// be created with New (or Async / AsyncWithContext) — the zero value of
// Future is not usable.
func New[T any]() *Future[T] {
	return &Future[T]{
		done: make(chan struct{}),
	}
}

// AsyncWithContext starts provider in a new goroutine and returns a Future that
// will hold the result once the goroutine completes. The provided context is
// forwarded to provider, allowing the async work to respect cancellation and
// deadlines set by the caller.
func AsyncWithContext[T any](ctx context.Context, provider func(context.Context) (T, error)) *Future[T] {
	f := New[T]()
	go func() {
		f.Resolve(provider(ctx))
	}()
	return f
}

// Async starts provider in a new goroutine and returns a Future that will
// hold the result once the goroutine completes.
func Async[T any](provider func() (T, error)) *Future[T] {
	f := New[T]()
	go func() {
		f.Resolve(provider())
	}()
	return f
}

// AwaitWithContext blocks until the async computation finishes and returns its result,
// or until ctx is cancelled — whichever happens first.
//
// Multiple goroutines may call AwaitWithContext on the same Future concurrently. All are
// unblocked when the result is ready and all receive the same value and error.
//
// If ctx is cancelled before the computation completes, AwaitWithContext returns the zero
// value of T together with ctx.Err(). The Future remains live: a subsequent
// call with a fresh context will still receive the result once it is available.
func (f *Future[T]) AwaitWithContext(ctx context.Context) (T, error) {
	select {
	case <-f.done:
		return f.val, f.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// AwaitWithTimeout blocks until the async computation finishes and returns its
// result, or until the given duration elapses — whichever happens first.
func (f *Future[T]) AwaitWithTimeout(d time.Duration) (T, error) {
	t := time.NewTimer(d)
	defer t.Stop()

	select {
	case <-f.done:
		return f.val, f.err
	case <-t.C:
		var zero T
		return zero, context.DeadlineExceeded
	}
}

// Await blocks until the async computation finishes and returns its result.
// It is equivalent to calling AwaitWithContext with context.Background().
func (f *Future[T]) Await() (T, error) {
	<-f.done
	return f.val, f.err
}

// Done reports whether the async computation has completed.
func (f *Future[T]) Done() bool {
	select {
	case <-f.done:
		return true
	default:
		return false
	}
}

// Resolve completes the Future with the given value and error, unblocking
// all callers blocked in Await, AwaitWithContext, or AwaitWithTimeout.
// Only the first call has any effect; subsequent calls are no-ops.
func (f *Future[T]) Resolve(val T, err error) {
	f.once.Do(func() {
		f.val, f.err = val, err
		close(f.done)
	})
}
