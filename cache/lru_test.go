package cache

import (
	"errors"
	"sync"
	"testing"
)

// mustNewLRUCache is a test helper that creates an LRUCache or fails the test.
func mustNewLRUCache[K comparable, V any](t testing.TB, capacity int) *LRUCache[K, V] {
	t.Helper()
	c, err := NewLRUCache[K, V](capacity)
	if err != nil {
		t.Fatalf("NewLRUCache(%d): unexpected error: %v", capacity, err)
	}
	return c
}

func TestNewLRUCache(t *testing.T) {
	t.Run("valid capacity", func(t *testing.T) {
		// when
		result := mustNewLRUCache[string, int](t, 3)
		// then
		if result.Cap() != 3 {
			t.Fatalf("Expected capacity 3, got %d", result.Cap())
		}

		if result.Len() != 0 {
			t.Fatalf("Expected empty cache, got len %d", result.Len())
		}
	})

	t.Run("invalid capacity", func(t *testing.T) {
		for _, cap := range []int{0, -1, -100} {
			_, err := NewLRUCache[string, int](cap)
			if !errors.Is(err, ErrInvalidCapacity) {
				t.Fatalf("capacity %d: expected ErrInvalidCapacity, got %v", cap, err)
			}
		}
	})
}

func TestLRUCachePut(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 3)
	// when
	c.Put("a", 1)
	c.Put("b", 2)
	// then
	if c.Len() != 2 {
		t.Fatalf("Expected len 2, got %d", c.Len())
	}
}

func TestLRUCacheGet(t *testing.T) {
	t.Run("existing key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		c.Put("a", 1)
		// when
		result, ok := c.Get("a")
		// then
		if !ok {
			t.Fatalf("Get returned false for existing key")
		}

		if result != 1 {
			t.Fatalf("Expected 1, got %d", result)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		// when
		result, ok := c.Get("missing")
		// then
		if ok {
			t.Fatalf("Get should return false for missing key")
		}

		if result != 0 {
			t.Fatalf("Get miss should return zero value, got %d", result)
		}
	})
}

func TestLRUCacheEvictsLRU(t *testing.T) {
	// given: capacity 2, insert a then b
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: insert c, which triggers eviction of a (LRU)
	c.Put("c", 3)
	// then
	if c.Len() != 2 {
		t.Fatalf("Expected len 2 after eviction, got %d", c.Len())
	}

	if c.Exists("a") {
		t.Fatalf("Key 'a' should have been evicted")
	}

	if !c.Exists("b") {
		t.Fatalf("Key 'b' should still be present")
	}

	if !c.Exists("c") {
		t.Fatalf("Key 'c' should be present")
	}
}

func TestLRUCacheGetPromotesToHead(t *testing.T) {
	// given: capacity 2, insert a then b (a is LRU)
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: access a (promotes it), then insert c
	c.Get("a")
	c.Put("c", 3)
	// then: b should be evicted (it is now LRU), a and c should survive
	if c.Exists("b") {
		t.Fatalf("Key 'b' should have been evicted (LRU after 'a' was promoted)")
	}

	if !c.Exists("a") {
		t.Fatalf("Key 'a' should still be present")
	}

	if !c.Exists("c") {
		t.Fatalf("Key 'c' should be present")
	}
}

func TestLRUCachePutUpdatePromotesToHead(t *testing.T) {
	// given: capacity 2, insert a then b (a is LRU)
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: update a (should promote it, not evict it)
	c.Put("a", 99)
	c.Put("c", 3)
	// then: b is evicted, a (with updated value) and c remain
	if c.Exists("b") {
		t.Fatalf("Key 'b' should have been evicted")
	}

	result, ok := c.Get("a")
	if !ok {
		t.Fatalf("Key 'a' should still be present after update")
	}

	if result != 99 {
		t.Fatalf("Expected updated value 99, got %d", result)
	}
}

func TestLRUCachePutUpdateDoesNotEvict(t *testing.T) {
	// given: cache at capacity
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: update an existing key (should not trigger eviction)
	c.Put("b", 20)
	// then: len stays at 2, no entries evicted
	if c.Len() != 2 {
		t.Fatalf("Updating existing key should not increase len, got %d", c.Len())
	}

	if !c.Exists("a") {
		t.Fatalf("Key 'a' should not have been evicted on update")
	}
}

func TestLRUCacheGetHeadTwice(t *testing.T) {
	// Regression: accessing the MRU node twice must not corrupt backward links.
	// given: three entries, c is MRU
	c := mustNewLRUCache[string, int](t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	c.Get("c") // c is now MRU
	// when: get c again (already at head)
	c.Get("c")
	// then: all entries accessible and eviction still correct
	if c.Len() != 3 {
		t.Fatalf("Expected len 3, got %d", c.Len())
	}

	c.Put("d", 4) // should evict a (LRU)
	if c.Exists("a") {
		t.Fatalf("Key 'a' should have been evicted")
	}

	if !c.Exists("b") {
		t.Fatalf("Key 'b' should still be present")
	}

	if !c.Exists("c") {
		t.Fatalf("Key 'c' should still be present")
	}

	if !c.Exists("d") {
		t.Fatalf("Key 'd' should be present")
	}
}

func TestLRUCacheCapacityOne(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 1)
	c.Put("a", 1)
	// when
	c.Get("a") // access the only node (regression: must not set tail=nil)
	c.Put("b", 2)
	// then
	if c.Exists("a") {
		t.Fatalf("Key 'a' should have been evicted")
	}

	result, ok := c.Get("b")
	if !ok {
		t.Fatalf("Key 'b' should be present")
	}

	if result != 2 {
		t.Fatalf("Expected 2, got %d", result)
	}

	if c.Len() != 1 {
		t.Fatalf("Expected len 1, got %d", c.Len())
	}
}

func TestLRUCacheDelete(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	// when
	c.Delete("a")
	// then
	if c.Exists("a") {
		t.Fatalf("Key 'a' should be gone after Delete")
	}

	if c.Len() != 1 {
		t.Fatalf("Expected len 1 after delete, got %d", c.Len())
	}

	_, ok := c.Get("a")
	if ok {
		t.Fatalf("Get should return false for deleted key")
	}
}

func TestLRUCacheDeleteMiss(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 3)
	c.Put("a", 1)
	// when: delete a key that doesn't exist (should not panic or corrupt state)
	c.Delete("missing")
	// then
	if c.Len() != 1 {
		t.Fatalf("Delete of missing key should not change len, got %d", c.Len())
	}

	if !c.Exists("a") {
		t.Fatalf("Key 'a' should still be present")
	}
}

func TestLRUCacheExists(t *testing.T) {
	t.Run("present key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		c.Put("a", 1)
		// when
		result := c.Exists("a")
		// then
		if !result {
			t.Fatalf("Exists should return true for present key")
		}
	})

	t.Run("absent key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		// when
		result := c.Exists("missing")
		// then
		if result {
			t.Fatalf("Exists should return false for absent key")
		}
	})
}

func TestLRUCacheExistsDoesNotPromote(t *testing.T) {
	// given: capacity 2, a is LRU
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: Exists on a (should not promote it)
	c.Exists("a")
	c.Put("c", 3) // triggers eviction of LRU
	// then: a is still evicted (Exists did not promote it)
	if c.Exists("a") {
		t.Fatalf("Exists should not promote 'a'; it should still be evicted")
	}
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	// given
	const goroutines = 100
	c := mustNewLRUCache[int, int](t, 50)
	var wg sync.WaitGroup
	// when: concurrent puts and gets
	for i := range goroutines {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			c.Put(v, v*10)
			c.Get(v)
		}(i)
	}
	wg.Wait()
	// then: cache len must be within capacity and not corrupt
	if c.Len() > c.Cap() {
		t.Fatalf("Cache len %d exceeds capacity %d", c.Len(), c.Cap())
	}
}
