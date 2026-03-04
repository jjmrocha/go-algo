package cache

import (
	"sync"

	"github.com/jjmrocha/go-algo/future"
)

// computations deduplicates concurrent provider calls for the same key.
// Each in-flight computation is represented by a [future.Future]; concurrent
// callers for the same key receive the same Future and block until it resolves.
// Completed computations are evicted from the map automatically.
type computations[K comparable, V any] struct {
	cache map[K]*future.Future[V]
	mu    sync.Mutex
}

// newComputations creates an empty computations registry.
func newComputations[K comparable, V any]() *computations[K, V] {
	return &computations[K, V]{
		cache: make(map[K]*future.Future[V]),
	}
}

// compute returns the in-flight Future for key if one exists, or starts a new
// background computation via provider and registers its Future. The Future is
// removed from the registry automatically when the computation completes.
func (c *computations[K, V]) compute(
	key K,
	provider func() (V, error),
) *future.Future[V] {
	c.mu.Lock()
	defer c.mu.Unlock()

	if f, ok := c.cache[key]; ok {
		return f
	}

	fn := func() (V, error) {
		v, err := provider()

		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.cache, key)

		return v, err
	}

	f := future.Async(fn)

	c.cache[key] = f

	return f
}
