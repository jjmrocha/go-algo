// Package cache provides in-memory caching implementations.
package cache

import (
	"sync"
	"time"
)

type node[K comparable, V any] struct {
	key    K
	value  V
	expire time.Time
	prev   *node[K, V]
	next   *node[K, V]
}

func (n *node[K, V]) expired(now time.Time) bool {
	if n.expire.IsZero() {
		return false
	}

	return now.After(n.expire)
}

// LRUCache is a thread-safe, generic Least Recently Used (LRU) cache.
// It evicts the least recently accessed entry when cap is exceeded.
// All load and write operations are O(1) average time complexity.
//
// K must be a comparable type; V can be any type.
type LRUCache[K comparable, V any] struct {
	cap  int
	ttl  time.Duration
	mu   sync.RWMutex
	kv   map[K]*node[K, V]
	head *node[K, V]
	tail *node[K, V]
	now  func() time.Time // injectable clock; defaults to time.Now
}

// NewLRUCache creates a new LRUCache configured by the supplied options.
// [WithCapacity] is required; [WithTTL] is optional.
//
// Returns [ErrInvalidCapacity] if [WithCapacity] is not provided or its value
// is less than or equal to zero. Returns [ErrInvalidTTL] if [WithTTL] is
// provided with a value less than or equal to zero.
func NewLRUCache[K comparable, V any](opts ...Option) (*LRUCache[K, V], error) {
	cfg, err := applyOptions(opts)
	if err != nil {
		return nil, err
	}

	return &LRUCache[K, V]{
		cap: cfg.capacity,
		ttl: cfg.ttl,
		kv:  make(map[K]*node[K, V]),
		now: time.Now,
	}, nil
}

// Get retrieves the value for key and marks it as most recently used.
// Returns the value and true if found, or the zero value and false if not.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.load(key)
}

// Put inserts or updates the value for key, marking it as most recently used.
// If the cache is at capacity and a new key is inserted, the least recently used entry is evicted.
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.write(key, value)
}

// Delete removes the entry for key. If key does not exist, it is a no-op.
func (c *LRUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.remove(key)
}

// Len returns the number of entries currently in the cache.
// When TTL is configured, this count may include recently-expired entries that
// have not yet been lazily evicted by a Get call.
func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.kv)
}

// Exists reports whether key is present in the cache without affecting LRU order.
func (c *LRUCache[K, V]) Exists(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	node, exists := c.kv[key]
	return exists && !node.expired(c.now())
}

// Cap returns the maximum number of entries the cache can hold.
func (c *LRUCache[K, V]) Cap() int {
	return c.cap
}

func (c *LRUCache[K, V]) load(key K) (V, bool) {
	n := c.kv[key]
	if n == nil {
		var zero V
		return zero, false
	}

	if n.expired(c.now()) {
		c.remove(key)
		var zero V
		return zero, false
	}

	c.moveToHead(n)
	return n.value, true
}

func (c *LRUCache[K, V]) moveToHead(n *node[K, V]) {
	if c.head == n {
		// Node is already at the head, no need to move
		return
	}

	// Remove n from its current position
	prev := n.prev
	next := n.next

	if prev != nil {
		prev.next = next
	}

	if next != nil {
		next.prev = prev
	}

	// Insert n at the head
	second := c.head
	n.next = second
	n.prev = nil
	second.prev = n
	c.head = n

	// Update tail if necessary
	if c.tail == n {
		c.tail = prev
	}
}

func (c *LRUCache[K, V]) write(key K, value V) {
	if node, exists := c.kv[key]; exists {
		node.value = value
		node.expire = c.expiresAt()
		c.moveToHead(node)
		return
	}

	// Create a new node and add it to the head
	newNode := &node[K, V]{
		key:    key,
		value:  value,
		expire: c.expiresAt(),
	}
	c.kv[key] = newNode

	if c.head == nil {
		c.head = newNode
		c.tail = newNode
	} else {
		newNode.next = c.head
		c.head.prev = newNode
		c.head = newNode
	}

	// Evict the least recently used item if cap is exceeded
	if len(c.kv) > c.cap {
		c.remove(c.tail.key)
	}
}

func (c *LRUCache[K, V]) remove(key K) {
	n, exists := c.kv[key]
	if !exists {
		// Key not found, nothing to remove
		return
	}

	// Remove n from the linked list
	prev := n.prev
	next := n.next

	if prev != nil {
		prev.next = next
	}

	if next != nil {
		next.prev = prev
	}

	// Update head and tail if necessary
	if c.head == n {
		c.head = next
	}

	if c.tail == n {
		c.tail = prev
	}

	// Remove the node from the kv
	delete(c.kv, key)
}

func (c *LRUCache[K, V]) expiresAt() time.Time {
	if c.ttl == 0 {
		return time.Time{}
	}

	return c.now().Add(c.ttl)
}
