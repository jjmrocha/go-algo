package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

// mustNewLRUWithProvider is a test helper that creates an LRUProvider or fails the test.
func mustNewLRUWithProvider[K comparable, V any](t testing.TB, capacity int, p Provider[K, V]) *LRUProvider[K, V] {
	t.Helper()
	lp, err := NewLRUWithProvider[K, V](capacity, p)
	if err != nil {
		t.Fatalf("NewLRUWithProvider(%d): unexpected error: %v", capacity, err)
	}
	return lp
}

// counterProvider returns a Provider that increments a call counter and
// returns the key as value (for string→string) or a supplied static value.
func counterProvider[K comparable](calls *atomic.Int32) Provider[K, K] {
	return func(key K) (K, error) {
		calls.Add(1)
		return key, nil
	}
}

// ---------- constructor ----------

func TestNewLRUWithProvider(t *testing.T) {
	t.Run("nil provider", func(t *testing.T) {
		// when
		_, result := NewLRUWithProvider[string, int](3, nil)
		// then
		if !errors.Is(result, ErrNilProvider) {
			t.Fatalf("expected ErrNilProvider, got %v", result)
		}
	})

	t.Run("invalid capacity", func(t *testing.T) {
		p := func(_ string) (int, error) { return 0, nil }
		for _, cap := range []int{0, -1} {
			_, err := NewLRUWithProvider[string, int](cap, p)
			if !errors.Is(err, ErrInvalidCapacity) {
				t.Fatalf("capacity %d: expected ErrInvalidCapacity, got %v", cap, err)
			}
		}
	})

	t.Run("cap and len", func(t *testing.T) {
		// given
		var calls atomic.Int32
		// when
		result := mustNewLRUWithProvider(t, 5, counterProvider[string](&calls))
		// then
		if result.Cap() != 5 {
			t.Fatalf("expected Cap 5, got %d", result.Cap())
		}

		if result.Len() != 0 {
			t.Fatalf("expected Len 0, got %d", result.Len())
		}
	})
}

// ---------- Get ----------

func TestLRUProviderGet(t *testing.T) {
	t.Run("cache miss", func(t *testing.T) {
		// given
		var calls atomic.Int32
		lp := mustNewLRUWithProvider(t, 3, counterProvider[string](&calls))
		// when
		result, err := lp.Get(context.Background(), "x")
		// then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != "x" {
			t.Fatalf("expected 'x', got %q", result)
		}

		if calls.Load() != 1 {
			t.Fatalf("expected 1 provider call, got %d", calls.Load())
		}
	})

	t.Run("cache hit", func(t *testing.T) {
		// given: first call populates cache
		var calls atomic.Int32
		lp := mustNewLRUWithProvider(t, 3, counterProvider[string](&calls))
		lp.Get(context.Background(), "x") //nolint:errcheck
		// when: second call should hit cache
		result, err := lp.Get(context.Background(), "x")
		// then
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != "x" {
			t.Fatalf("expected 'x', got %q", result)
		}

		if calls.Load() != 1 {
			t.Fatalf("expected exactly 1 provider call (cache hit on second), got %d", calls.Load())
		}
	})

	t.Run("provider error", func(t *testing.T) {
		// given
		providerErr := errors.New("provider failed")
		p := func(_ string) (int, error) { return 0, providerErr }
		lp := mustNewLRUWithProvider(t, 3, p)
		// when
		result, err := lp.Get(context.Background(), "x")
		// then
		if !errors.Is(err, providerErr) {
			t.Fatalf("expected providerErr, got %v", err)
		}

		if result != 0 {
			t.Fatalf("expected zero value on error, got %d", result)
		}
	})

	t.Run("error does not cache", func(t *testing.T) {
		// given: provider fails on first call, succeeds on second
		var calls atomic.Int32
		providerErr := errors.New("transient error")
		p := func(k string) (string, error) {
			n := calls.Add(1)
			if n == 1 {
				return "", providerErr
			}
			return k, nil
		}
		lp := mustNewLRUWithProvider(t, 3, p)
		// when: first call fails
		_, err := lp.Get(context.Background(), "x")
		if !errors.Is(err, providerErr) {
			t.Fatalf("expected providerErr on first call, got %v", err)
		}
		// then: second call should retry the provider (not return a cached error)
		result, err := lp.Get(context.Background(), "x")
		if err != nil {
			t.Fatalf("unexpected error on second call: %v", err)
		}

		if result != "x" {
			t.Fatalf("expected 'x', got %q", result)
		}

		if calls.Load() != 2 {
			t.Fatalf("expected 2 provider calls, got %d", calls.Load())
		}
	})
}

// ---------- Invalidate ----------

func TestLRUProviderInvalidate(t *testing.T) {
	t.Run("forces provider re-call", func(t *testing.T) {
		// given: populate cache
		var calls atomic.Int32
		lp := mustNewLRUWithProvider(t, 3, counterProvider[string](&calls))
		lp.Get(context.Background(), "x") //nolint:errcheck
		// when: invalidate and get again
		lp.Invalidate("x")
		lp.Get(context.Background(), "x") //nolint:errcheck
		// then: provider must have been called twice
		if calls.Load() != 2 {
			t.Fatalf("expected 2 provider calls after invalidation, got %d", calls.Load())
		}
	})

	t.Run("missing key is no-op", func(t *testing.T) {
		// given
		var calls atomic.Int32
		lp := mustNewLRUWithProvider(t, 3, counterProvider[string](&calls))
		// when: invalidate a key that was never fetched (should be a no-op)
		lp.Invalidate("missing")
		// then: len and state unaffected
		if lp.Len() != 0 {
			t.Fatalf("expected Len 0 after no-op Invalidate, got %d", lp.Len())
		}
	})
}

// ---------- Len / Cap ----------

func TestLRUProviderLen(t *testing.T) {
	// given
	lp := mustNewLRUWithProvider(t, 3, counterProvider[string](new(atomic.Int32)))
	// when
	lp.Get(context.Background(), "a") //nolint:errcheck
	lp.Get(context.Background(), "b") //nolint:errcheck
	// then
	if lp.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", lp.Len())
	}
}

func TestLRUProviderCap(t *testing.T) {
	// when
	result := mustNewLRUWithProvider(t, 7, counterProvider[string](new(atomic.Int32)))
	// then
	if result.Cap() != 7 {
		t.Fatalf("expected Cap 7, got %d", result.Cap())
	}
}

// ---------- Eviction ----------

func TestLRUProviderEviction(t *testing.T) {
	// given: capacity 2, fill with a and b
	var calls atomic.Int32
	lp := mustNewLRUWithProvider(t, 2, counterProvider[string](&calls))
	lp.Get(context.Background(), "a") //nolint:errcheck
	lp.Get(context.Background(), "b") //nolint:errcheck
	// when: fetch c — evicts a (LRU)
	lp.Get(context.Background(), "c") //nolint:errcheck
	// then: fetching a again must call the provider (it was evicted)
	callsBefore := calls.Load()
	lp.Get(context.Background(), "a") //nolint:errcheck
	if calls.Load() == callsBefore {
		t.Fatalf("expected provider call for evicted key 'a'")
	}

	if lp.Len() > lp.Cap() {
		t.Fatalf("Len %d exceeds Cap %d", lp.Len(), lp.Cap())
	}
}

// ---------- Concurrency ----------

func TestLRUProviderConcurrentGet(t *testing.T) {
	// given
	const goroutines = 200
	var calls atomic.Int32
	lp := mustNewLRUWithProvider(t, 50, counterProvider[int](&calls))
	var wg sync.WaitGroup
	// when: concurrent Gets on overlapping keys
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			lp.Get(context.Background(), v%20) //nolint:errcheck
		}(i)
	}
	wg.Wait()
	// then: cache is not corrupted and len is within capacity
	if lp.Len() > lp.Cap() {
		t.Fatalf("Len %d exceeds Cap %d after concurrent access", lp.Len(), lp.Cap())
	}
}
