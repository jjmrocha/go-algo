// Package cache provides in-memory caching implementations.
package cache

import "sync"

type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// LRUCache is a thread-safe, generic Least Recently Used (LRU) cache.
// It evicts the least recently accessed entry when capacity is exceeded.
// All read and write operations are O(1) average time complexity.
//
// K must be a comparable type; V can be any type.
type LRUCache[K comparable, V any] struct {
	capacity int
	mu       sync.RWMutex
	cache    map[K]*node[K, V]
	head     *node[K, V]
	tail     *node[K, V]
}

// NewLRUCache creates a new LRUCache with the given capacity.
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*node[K, V]),
	}
}

// Get retrieves the value for key and marks it as most recently used.
// Returns the value and true if found, or the zero value and false if not.
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.read(key)
}

// Put inserts or updates the value for key, marking it as most recently used.
// If the cache is at capacity and a new key is inserted, the least recently used entry is evicted.
func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store(key, value)
}

// Delete removes the entry for key. If key does not exist, it is a no-op.
func (c *LRUCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.remove(key)
}

// Len returns the number of entries currently in the cache.
func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// Exists reports whether key is present in the cache without affecting LRU order.
func (c *LRUCache[K, V]) Exists(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.cache[key]
	return exists
}

// Cap returns the maximum number of entries the cache can hold.
func (c *LRUCache[K, V]) Cap() int {
	return c.capacity
}

func (c *LRUCache[K, V]) read(key K) (V, bool) {
	node := c.cache[key]
	if node == nil {
		var zero V
		return zero, false
	}

	c.moveToHead(node)
	return node.value, true
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

func (c *LRUCache[K, V]) store(key K, value V) {
	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToHead(node)
		return
	}

	// Create a new node and add it to the head
	newNode := &node[K, V]{key: key, value: value}
	c.cache[key] = newNode

	if c.head == nil {
		c.head = newNode
		c.tail = newNode
	} else {
		newNode.next = c.head
		c.head.prev = newNode
		c.head = newNode
	}

	// Evict the least recently used item if capacity is exceeded
	if len(c.cache) > c.capacity {
		c.remove(c.tail.key)
	}
}

func (c *LRUCache[K, V]) remove(key K) {
	n, exists := c.cache[key]
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

	// Remove the node from the cache
	delete(c.cache, key)
}
