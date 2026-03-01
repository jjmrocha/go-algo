package cache

import (
	"context"
	"errors"
)

// Provider is a function that fetches the value for key on a cache miss.
// It returns the fetched value and nil on success, or the zero value and a
// non-nil error on failure. Errors are not cached; a failed fetch will be
// retried on the next [LRUProvider.Get] call for the same key.
type Provider[K comparable, V any] func(key K) (V, error)

// LRUProvider wraps an [LRUCache] with a lazy-loading [Provider], implementing
// the cache-aside pattern. On a cache miss the provider is called in a
// background goroutine to fetch the value, which is then stored in the cache
// for subsequent reads.
//
// Concurrent [LRUProvider.Get] calls for the same key are deduplicated: only
// one provider invocation runs at a time and its result is shared with all
// callers waiting on the same key.
//
// LRUProvider is safe for concurrent use by multiple goroutines.
type LRUProvider[K comparable, V any] struct {
	cache    *LRUCache[K, V]
	fetches  *computations[K, V]
	provider Provider[K, V]
}

// ErrNilProvider is returned by [NewLRUWithProvider] when a nil provider is supplied.
var ErrNilProvider = errors.New("provider cannot be nil")

// NewLRUWithProvider creates an [LRUProvider] backed by an LRU cache of the
// given capacity and the supplied provider function.
//
// Returns [ErrInvalidCapacity] if capacity is less than or equal to zero, or
// [ErrNilProvider] if p is nil.
func NewLRUWithProvider[K comparable, V any](capacity int, p Provider[K, V]) (*LRUProvider[K, V], error) {
	if p == nil {
		return nil, ErrNilProvider
	}

	lru, err := NewLRUCache[K, V](capacity)
	if err != nil {
		return nil, err
	}

	return &LRUProvider[K, V]{
		cache:    lru,
		provider: p,
		fetches:  newComputations[K, V](),
	}, nil
}

// Get returns the cached value for key. On a cache miss the provider is called
// in a background goroutine to fetch the value, which is stored in the cache
// before being returned to all concurrent callers waiting on the same key.
//
// ctx controls how long the caller is willing to wait for the result. If ctx
// is cancelled before the computation completes, Get returns the zero value of
// V and ctx.Err(). The background computation continues and its result will be
// available to subsequent calls.
//
// Provider errors are not cached; a subsequent Get will retry the provider.
func (lp *LRUProvider[K, V]) Get(ctx context.Context, key K) (V, error) {
	if value, ok := lp.cache.Get(key); ok {
		return value, nil
	}

	return lp.fetches.compute(key, func() (V, error) {
		value, err := lp.provider(key)
		if err != nil {
			var zero V
			return zero, err
		}

		lp.cache.Put(key, value)
		return value, nil
	}).Await(ctx)
}

// Invalidate removes key from the LRU cache. The next [LRUProvider.Get] for
// that key will invoke the provider. If key is not present, Invalidate is a
// no-op.
//
// Note: if a background computation for key is already in flight when
// Invalidate is called, the computed value will still be inserted into the
// cache once the computation completes.
func (lp *LRUProvider[K, V]) Invalidate(key K) {
	lp.cache.Delete(key)
}

// Len returns the number of entries currently held in the cache.
func (lp *LRUProvider[K, V]) Len() int {
	return lp.cache.Len()
}

// Cap returns the maximum number of entries the cache can hold.
func (lp *LRUProvider[K, V]) Cap() int {
	return lp.cache.Cap()
}
