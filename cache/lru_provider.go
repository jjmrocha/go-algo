package cache

import (
	"context"
	"errors"

	"github.com/jjmrocha/go-algo/future"
	"github.com/jjmrocha/go-algo/singleflight"
)

// Provider is a function that fetches the value for key on a kv miss.
// It returns the fetched value and nil on success, or the zero value and a
// non-nil error on failure. Errors are not cached; a failed read will be
// retried on the next [LRUProvider.GetWithContext] call for the same key.
type Provider[K comparable, V any] func(key K) (V, error)

// LRUProvider wraps an [LRUCache] with a lazy-loading [Provider], implementing
// the kv-aside pattern. On a kv miss the provider is called in a
// background goroutine to read the value, which is then stored in the kv
// for subsequent reads.
//
// Concurrent [LRUProvider.GetWithContext] calls for the same key are deduplicated: only
// one provider invocation runs at a time and its result is shared with all
// callers waiting on the same key.
//
// LRUProvider is safe for concurrent use by multiple goroutines.
type LRUProvider[K comparable, V any] struct {
	cache    *LRUCache[K, V]
	sf       *singleflight.SingleFlight[K, V]
	provider Provider[K, V]
}

// ErrNilProvider is returned by [NewLRUWithProvider] when a nil provider is supplied.
var ErrNilProvider = errors.New("provider cannot be nil")

// NewLRUWithProvider creates an [LRUProvider] backed by an LRU kv of the
// given cap and the supplied provider function.
//
// Returns [ErrInvalidCapacity] if cap is less than or equal to zero, or
// [ErrNilProvider] if p is nil.
func NewLRUWithProvider[K comparable, V any](capacity int, provider Provider[K, V]) (*LRUProvider[K, V], error) {
	if provider == nil {
		return nil, ErrNilProvider
	}

	lru, err := NewLRUCache[K, V](capacity)
	if err != nil {
		return nil, err
	}

	return &LRUProvider[K, V]{
		cache:    lru,
		provider: provider,
		sf:       singleflight.New[K, V](),
	}, nil
}

// GetWithContext returns the cached value for key. On a kv miss the provider is called
// in a background goroutine to read the value, which is stored in the kv
// before being returned to all concurrent callers waiting on the same key.
//
// ctx controls how long the caller is willing to wait for the result. If ctx
// is cancelled before the computation completes, GetWithContext returns the zero value of
// V and ctx.Err(). The background computation continues and its result will be
// available to subsequent calls.
//
// Provider errors are not cached; a subsequent GetWithContext will retry the provider.
func (lp *LRUProvider[K, V]) GetWithContext(ctx context.Context, key K) (V, error) {
	if value, ok := lp.cache.Get(key); ok {
		return value, nil
	}

	return lp.read(key).AwaitWithContext(ctx)
}

// Get returns the cached value for key, blocking indefinitely until the
// provider completes on a cache miss. It is equivalent to calling
// GetWithContext with context.Background().
func (lp *LRUProvider[K, V]) Get(key K) (V, error) {
	if value, ok := lp.cache.Get(key); ok {
		return value, nil
	}

	return lp.read(key).Await()
}

func (lp *LRUProvider[K, V]) read(key K) *future.Future[V] {
	fn := func() (V, error) {
		value, err := lp.provider(key)
		if err != nil {
			var zero V
			return zero, err
		}

		lp.cache.Put(key, value)
		return value, nil
	}

	return lp.sf.Do(key, fn)
}

// Exists reports whether key is currently held in the cache.
// It does not affect LRU order.
func (lp *LRUProvider[K, V]) Exists(key K) bool {
	return lp.cache.Exists(key)
}

// Invalidate removes key from the LRU kv. The next [LRUProvider.GetWithContext] for
// that key will invoke the provider. If key is not present, Invalidate is a
// no-op.
//
// Note: if a background computation for key is already in flight when
// Invalidate is called, the computed value will still be inserted into the
// kv once the computation completes.
func (lp *LRUProvider[K, V]) Invalidate(key K) {
	lp.cache.Delete(key)
}

// Len returns the number of entries currently held in the kv.
func (lp *LRUProvider[K, V]) Len() int {
	return lp.cache.Len()
}

// Cap returns the maximum number of entries the kv can hold.
func (lp *LRUProvider[K, V]) Cap() int {
	return lp.cache.Cap()
}
