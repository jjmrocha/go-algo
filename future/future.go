// Package future provides a generic Future/Promise implementation for Go,
// enabling asynchronous computation with context-aware cancellation.
//
// A Future represents a value that will be available at some point in the
// future, produced by an asynchronous operation. The caller blocks until the
// result is ready via [Future.Await].
//
// Example usage:
//
//	f := future.AsyncWithContext(ctx, func(ctx context.Context) (int, error) {
//	    return expensiveComputation(ctx)
//	})
//
//	result, err := f.Await(ctx)
package future

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// ErrAwaitAlreadyCalled is returned when the Await method was already called
// for the future
var ErrAwaitAlreadyCalled = errors.New("await already called")

// payload is the internal carrier for a computation result, bundling the
// produced value and any error returned by the async provider.
type payload[T any] struct {
	val T
	err error
}

// Future represents the result of an asynchronous computation of type T.
// Concurrent calls to [Future.Await] are safe: the first caller receives the
// result, subsequent callers get [ErrAwaitAlreadyCalled]. For fan-out
// scenarios, each consumer should obtain its own Future via [AsyncWithContext].
type Future[T any] struct {
	ch     chan payload[T]
	called atomic.Bool
}

// newFuture creates an uninitialised Future with a buffered channel of capacity 1,
// ensuring the producer goroutine never blocks even if Await returns early.
func newFuture[T any]() *Future[T] {
	return &Future[T]{
		ch: make(chan payload[T], 1),
	}
}

// AsyncWithContext starts provider in a new goroutine and returns a Future that will
// hold the result once the goroutine completes. The provided context is
// forwarded to provider, allowing the async work to respect cancellation
// deadlines set by the caller.
func AsyncWithContext[T any](ctx context.Context, provider func(context.Context) (T, error)) *Future[T] {
	f := newFuture[T]()
	go func() {
		f.resolve(provider(ctx))
	}()
	return f
}

// Async starts provider in a new goroutine and returns a Future that will
// hold the result once the goroutine completes.
func Async[T any](provider func() (T, error)) *Future[T] {
	f := newFuture[T]()
	go func() {
		f.resolve(provider())
	}()
	return f
}

func (f *Future[T]) resolve(val T, err error) {
	f.ch <- payload[T]{val, err}
	close(f.ch)
}

// Await blocks until the async computation finishes and returns its result,
// or until ctx is cancelled — whichever happens first.
//
// If ctx is cancelled before the computation completes, Await returns the
// zero value of T together with ctx.Err(). If the computation finishes first,
// Await returns the value and error produced by the provider passed to [AsyncWithContext].
func (f *Future[T]) Await(ctx context.Context) (T, error) {
	var zero T

	if f.called.Load() {
		return zero, ErrAwaitAlreadyCalled
	}

	select {
	case payload := <-f.ch:
		f.called.Store(true)
		return payload.val, payload.err
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

// AwaitWithTimeout blocks until the async computation finishes and returns its
// result, or until the given duration elapses — whichever happens first.
func (f *Future[T]) AwaitWithTimeout(t time.Duration) (T, error) {
	var zero T

	if f.called.Load() {
		return zero, ErrAwaitAlreadyCalled
	}

	select {
	case payload := <-f.ch:
		f.called.Store(true)
		return payload.val, payload.err
	case <-time.After(t):
		return zero, context.DeadlineExceeded
	}
}
