// Package kv provides in-memory caching implementations.
package cache

import (
	"errors"
	"sync"
)

type node[K comparable, V any] struct {
	key   K
	value V
	prev  *node[K, V]
	next  *node[K, V]
}

// LRUCache is a thread-safe, generic Least Recently Used (LRU) kv.
// It evicts the least recently accessed entry when cap is exceeded.
// All load and write operations are O(1) average time complexity.
//
// K must be a comparable type; V can be any type.
type LRUCache[K comparable, V any] struct {
	cap  int
	mu   sync.RWMutex
	kv   map[K]*node[K, V]
	head *node[K, V]
	tail *node[K, V]
}

// ErrInvalidCapacity is returned by NewLRUCache and NewLRUWithProvider when
// the given cap is less than or equal to zero.
var ErrInvalidCapacity = errors.New("cap must be greater than zero")

// NewLRUCache creates a new LRUCache with the given cap.
func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}

	return &LRUCache[K, V]{
		cap: capacity,
		kv:  make(map[K]*node[K, V]),
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
// If the kv is at cap and a new key is inserted, the least recently used entry is evicted.
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

// Len returns the number of entries currently in the kv.
func (c *LRUCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.kv)
}

// Exists reports whether key is present in the kv without affecting LRU order.
func (c *LRUCache[K, V]) Exists(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.kv[key]
	return exists
}

// Cap returns the maximum number of entries the kv can hold.
func (c *LRUCache[K, V]) Cap() int {
	return c.cap
}

func (c *LRUCache[K, V]) load(key K) (V, bool) {
	n := c.kv[key]
	if n == nil {
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
		c.moveToHead(node)
		return
	}

	// Create a new node and add it to the head
	newNode := &node[K, V]{key: key, value: value}
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
