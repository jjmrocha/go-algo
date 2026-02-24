// Package future provides a generic Future/Promise implementation for Go,
// enabling asynchronous computation with context-aware cancellation.
//
// A Future represents a value that will be available at some point in the
// future, produced by an asynchronous operation. The caller can either block
// until the result is ready ([Future.Await]) or register a callback to be
// invoked when the result becomes available ([Future.Then]).
//
// Example usage:
//
//	f := future.Async(ctx, func(ctx context.Context) (int, error) {
//	    return expensiveComputation(ctx)
//	})
//
//	result, err := f.Await(ctx)
package future

import "context"

// payload is the internal carrier for a computation result, bundling the
// produced value and any error returned by the async provider.
type payload[T any] struct {
	val T
	err error
}

// Future represents the result of an asynchronous computation of type T.
// It is safe to call [Future.Await] or [Future.Then] from multiple goroutines,
// but each call will race to consume the single result; only one caller will
// receive the value. For fan-out scenarios, each consumer should obtain its
// own Future via [Async].
type Future[T any] struct {
	ch chan payload[T]
}

// new creates an uninitialised Future with an unbuffered channel.
func new[T any]() *Future[T] {
	return &Future[T]{
		ch: make(chan payload[T]),
	}
}

// Async starts provider in a new goroutine and returns a Future that will
// hold the result once the goroutine completes. The provided context is
// forwarded to provider, allowing the async work to respect cancellation
// deadlines set by the caller.
//
// The returned Future must be consumed exactly once via [Future.Await] or
// [Future.Then] to avoid leaking the goroutine.
func Async[T any](ctx context.Context, provider func(context.Context) (T, error)) *Future[T] {
	f := new[T]()
	go func() {
		v, err := provider(ctx)
		f.ch <- payload[T]{v, err}
		close(f.ch)
	}()
	return f
}

// Await blocks until the async computation finishes and returns its result,
// or until ctx is cancelled — whichever happens first.
//
// If ctx is cancelled before the computation completes, Await returns the
// zero value of T together with ctx.Err(). If the computation finishes first,
// Await returns the value and error produced by the provider passed to [Async].
func (f *Future[T]) Await(ctx context.Context) (T, error) {
	select {
	case payload := <-f.ch:
		return payload.val, payload.err
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	}
}

// Then registers consumer as a callback to be invoked asynchronously when
// the Future's result is ready. The callback receives the same (value, error)
// pair that [Future.Await] would return, including any context cancellation
// error if ctx expires before the computation completes.
//
// Then returns immediately; consumer is called in a new goroutine.
func (f *Future[T]) Then(ctx context.Context, consumer func(T, error)) {
	go func() {
		consumer(f.Await(ctx))
	}()
}
