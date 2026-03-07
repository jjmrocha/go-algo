package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
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
	t.Run("valid cap", func(t *testing.T) {
		// when
		result := mustNewLRUCache[string, int](t, 3)
		// then
		assert.Equal(t, 3, result.Cap())
		assert.Equal(t, 0, result.Len())
	})

	t.Run("invalid cap", func(t *testing.T) {
		for _, cap := range []int{0, -1, -100} {
			// when
			_, err := NewLRUCache[string, int](cap)
			// then
			assert.ErrorIs(t, err, ErrInvalidCapacity, "cap %d", cap)
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
	assert.Equal(t, 2, c.Len())
}

func TestLRUCacheGet(t *testing.T) {
	t.Run("existing key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		c.Put("a", 1)
		// when
		result, ok := c.Get("a")
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, result)
	})

	t.Run("missing key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		// when
		result, ok := c.Get("missing")
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, result)
	})
}

func TestLRUCacheEvictsLRU(t *testing.T) {
	// given: cap 2, insert a then b
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: insert c, which triggers eviction of a (LRU)
	c.Put("c", 3)
	// then
	assert.Equal(t, 2, c.Len())
	assert.False(t, c.Exists("a"))
	assert.True(t, c.Exists("b"))
	assert.True(t, c.Exists("c"))
}

func TestLRUCacheGetPromotesToHead(t *testing.T) {
	// given: cap 2, insert a then b (a is LRU)
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: access a (promotes it), then insert c
	c.Get("a")
	c.Put("c", 3)
	// then: b should be evicted (it is now LRU), a and c should survive
	assert.False(t, c.Exists("b"))
	assert.True(t, c.Exists("a"))
	assert.True(t, c.Exists("c"))
}

func TestLRUCachePutUpdatePromotesToHead(t *testing.T) {
	// given: cap 2, insert a then b (a is LRU)
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: update a (should promote it, not evict it)
	c.Put("a", 99)
	c.Put("c", 3)
	// then: b is evicted, a (with updated value) and c remain
	assert.False(t, c.Exists("b"))
	result, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 99, result)
}

func TestLRUCachePutUpdateDoesNotEvict(t *testing.T) {
	// given: kv at cap
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: update an existing key (should not trigger eviction)
	c.Put("b", 20)
	// then: len stays at 2, no entries evicted
	assert.Equal(t, 2, c.Len())
	assert.True(t, c.Exists("a"))
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
	assert.Equal(t, 3, c.Len())
	c.Put("d", 4) // should evict a (LRU)
	assert.False(t, c.Exists("a"))
	assert.True(t, c.Exists("b"))
	assert.True(t, c.Exists("c"))
	assert.True(t, c.Exists("d"))
}

func TestLRUCacheCapacityOne(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 1)
	c.Put("a", 1)
	// when
	c.Get("a") // access the only node (regression: must not set tail=nil)
	c.Put("b", 2)
	// then
	assert.False(t, c.Exists("a"))
	result, ok := c.Get("b")
	assert.True(t, ok)
	assert.Equal(t, 2, result)
	assert.Equal(t, 1, c.Len())
}

func TestLRUCacheDelete(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 3)
	c.Put("a", 1)
	c.Put("b", 2)
	// when
	c.Delete("a")
	// then
	assert.False(t, c.Exists("a"))
	assert.Equal(t, 1, c.Len())
	_, ok := c.Get("a")
	assert.False(t, ok)
}

func TestLRUCacheDeleteMiss(t *testing.T) {
	// given
	c := mustNewLRUCache[string, int](t, 3)
	c.Put("a", 1)
	// when: delete a key that doesn't exist (should not panic or corrupt state)
	c.Delete("missing")
	// then
	assert.Equal(t, 1, c.Len())
	assert.True(t, c.Exists("a"))
}

func TestLRUCacheExists(t *testing.T) {
	t.Run("present key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		c.Put("a", 1)
		// when
		result := c.Exists("a")
		// then
		assert.True(t, result)
	})

	t.Run("absent key", func(t *testing.T) {
		// given
		c := mustNewLRUCache[string, int](t, 3)
		// when
		result := c.Exists("missing")
		// then
		assert.False(t, result)
	})
}

func TestLRUCacheExistsDoesNotPromote(t *testing.T) {
	// given: cap 2, a is LRU
	c := mustNewLRUCache[string, int](t, 2)
	c.Put("a", 1)
	c.Put("b", 2)
	// when: Exists on a (should not promote it)
	c.Exists("a")
	c.Put("c", 3) // triggers eviction of LRU
	// then: a is still evicted (Exists did not promote it)
	assert.False(t, c.Exists("a"))
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
	// then: kv len must be within cap and not corrupt
	assert.LessOrEqual(t, c.Len(), c.Cap())
}
