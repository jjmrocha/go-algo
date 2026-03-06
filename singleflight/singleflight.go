// Package singleflight deduplicates concurrent calls for the same key,
// preventing the cache stampede problem.
//
// When multiple goroutines concurrently request the same key, only one call to
// the underlying provider is made. All concurrent callers share the same
// [future.Future] and receive the same result when the computation completes.
//
// Unlike memoization, completed computations are not retained: a new call after
// completion starts a fresh provider invocation. SingleFlight is typically used
// alongside an external cache to prevent stampedes on cache misses.
package singleflight

import (
	"sync"

	"github.com/jjmrocha/go-algo/future"
)

// SingleFlight deduplicates concurrent provider calls for the same key.
// Each in-flight computation is represented by a [future.Future]; concurrent
// callers for the same key receive the same Future and share the result when
// the computation completes. Completed computations are evicted from the
// registry automatically.
//
// SingleFlight is safe for concurrent use by multiple goroutines.
type SingleFlight[K comparable, V any] struct {
	keys map[K]*future.Future[V]
	mu   sync.Mutex
}

// New creates an empty SingleFlight registry.
func New[K comparable, V any]() *SingleFlight[K, V] {
	return &SingleFlight[K, V]{
		keys: make(map[K]*future.Future[V]),
	}
}

// Do returns the in-flight [future.Future] for key if one exists, or starts a
// new background computation by calling provider in a goroutine and registers
// its Future. All concurrent callers for the same key receive the same Future
// and share the result.
//
// The entry is removed from the registry when the computation completes,
// regardless of whether provider returned an error. A subsequent call for the
// same key after completion starts a fresh computation.
func (s *SingleFlight[K, V]) Do(
	key K,
	provider func() (V, error),
) *future.Future[V] {
	s.mu.Lock()
	defer s.mu.Unlock()

	if f, ok := s.keys[key]; ok {
		return f
	}

	fn := func() (V, error) {
		v, err := provider()

		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.keys, key)

		return v, err
	}

	f := future.Async(fn)

	s.keys[key] = f

	return f
}
