package hashtable

import "golang.org/x/exp/constraints"

// SeparateChainingHashTable is a hash table implementation using separate chaining
type SeparateChainingHashTable[K constraints.Ordered, V any] struct {
	m    int                 // number of chains
	n    int                 // number of key-value pairs
	st   []*chain[K, V]      // array of chains
}

type chain[K constraints.Ordered, V any] struct {
	key   K
	value V
	next  *chain[K, V]
}

// NewSeparateChainingHashTable creates a new hash table with separate chaining
func NewSeparateChainingHashTable[K constraints.Ordered, V any](m int) *SeparateChainingHashTable[K, V] {
	if m < 1 {
		m = 997 // default prime number
	}
	return &SeparateChainingHashTable[K, V]{
		m:  m,
		st: make([]*chain[K, V], m),
	}
}

// Put inserts a key-value pair into the hash table
func (ht *SeparateChainingHashTable[K, V]) Put(key K, value V) {
	i := ht.hash(key)
	
	// Check if key already exists
	for x := ht.st[i]; x != nil; x = x.next {
		if x.key == key {
			x.value = value
			return
		}
	}
	
	// Insert at beginning of chain
	ht.st[i] = &chain[K, V]{key: key, value: value, next: ht.st[i]}
	ht.n++
}

// Get retrieves the value associated with the given key
func (ht *SeparateChainingHashTable[K, V]) Get(key K) (V, bool) {
	i := ht.hash(key)
	for x := ht.st[i]; x != nil; x = x.next {
		if x.key == key {
			return x.value, true
		}
	}
	var zero V
	return zero, false
}

// Delete removes a key and its associated value from the hash table
func (ht *SeparateChainingHashTable[K, V]) Delete(key K) {
	i := ht.hash(key)
	
	// Handle empty chain
	if ht.st[i] == nil {
		return
	}
	
	// Handle first node
	if ht.st[i].key == key {
		ht.st[i] = ht.st[i].next
		ht.n--
		return
	}
	
	// Search rest of chain
	for x := ht.st[i]; x.next != nil; x = x.next {
		if x.next.key == key {
			x.next = x.next.next
			ht.n--
			return
		}
	}
}

// Contains returns true if the hash table contains the given key
func (ht *SeparateChainingHashTable[K, V]) Contains(key K) bool {
	_, ok := ht.Get(key)
	return ok
}

// Size returns the number of key-value pairs in the hash table
func (ht *SeparateChainingHashTable[K, V]) Size() int {
	return ht.n
}

// IsEmpty returns true if the hash table is empty
func (ht *SeparateChainingHashTable[K, V]) IsEmpty() bool {
	return ht.n == 0
}

func (ht *SeparateChainingHashTable[K, V]) hash(key K) int {
	// Simple hash function for ordered types
	h := 0
	switch any(key).(type) {
	case int:
		h = int(any(key).(int))
	case string:
		s := any(key).(string)
		for i := 0; i < len(s); i++ {
			h = 31*h + int(s[i])
		}
	default:
		// For other types, use a basic conversion
		h = int(any(key).(int64))
	}
	
	// Make sure hash is positive
	if h < 0 {
		h = -h
	}
	return h % ht.m
}
