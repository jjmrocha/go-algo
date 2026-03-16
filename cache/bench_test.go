package cache

import (
	"fmt"
	"testing"
)

// ---------- LRUCache ----------

func BenchmarkLRUCacheGetHit(b *testing.B) {
	c, _ := NewLRUCache[string, int](WithCapacity(1024))
	c.Put("key", 42)
	for b.Loop() {
		c.Get("key")
	}
}

func BenchmarkLRUCacheGetMiss(b *testing.B) {
	c, _ := NewLRUCache[string, int](WithCapacity(1024))
	for b.Loop() {
		c.Get("missing")
	}
}

func BenchmarkLRUCachePutNew(b *testing.B) {
	c, _ := NewLRUCache[string, int](WithCapacity(1 << 20)) // large cap — no eviction
	for b.Loop() {
		c.Put("key", 1)
	}
}

func BenchmarkLRUCachePutEvict(b *testing.B) {
	// cap=1 so every Put of a new key evicts the previous one
	c, _ := NewLRUCache[string, int](WithCapacity(1))
	keys := [2]string{"a", "b"}
	i := 0
	for b.Loop() {
		c.Put(keys[i&1], i)
		i++
	}
}

func BenchmarkLRUCacheExists(b *testing.B) {
	c, _ := NewLRUCache[string, int](WithCapacity(1024))
	c.Put("key", 1)
	for b.Loop() {
		c.Exists("key")
	}
}

// ---------- LRUCache — parallel reads ----------

func BenchmarkLRUCacheGetParallel(b *testing.B) {
	c, _ := NewLRUCache[string, int](WithCapacity(1024))
	for i := range 64 {
		c.Put(fmt.Sprintf("k%d", i), i)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Get(fmt.Sprintf("k%d", i%64))
			i++
		}
	})
}

// ---------- LRUProvider ----------

func BenchmarkLRUProviderGetHit(b *testing.B) {
	calls := 0
	lp, _ := NewLRUWithProvider[string, int](func(key string) (int, error) {
		calls++
		return 42, nil
	}, WithCapacity(1024))
	_, _ = lp.Get("key") // warm the cache
	b.ResetTimer()
	for b.Loop() {
		_, _ = lp.Get("key")
	}
}

func BenchmarkLRUProviderGetMiss(b *testing.B) {
	lp, _ := NewLRUWithProvider[int, int](func(key int) (int, error) {
		return key, nil
	}, WithCapacity(1))
	i := 0
	for b.Loop() {
		lp.Get(i) //nolint:errcheck
		i++
	}
}
