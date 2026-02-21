package hashtable

import "golang.org/x/exp/constraints"

// LinearProbingHashTable is a hash table implementation using linear probing
type LinearProbingHashTable[K constraints.Ordered, V any] struct {
	m    int   // size of linear probing table
	n    int   // number of key-value pairs
	keys []K
	vals []V
}

// NewLinearProbingHashTable creates a new hash table with linear probing
func NewLinearProbingHashTable[K constraints.Ordered, V any](m int) *LinearProbingHashTable[K, V] {
	if m < 1 {
		m = 16
	}
	ht := &LinearProbingHashTable[K, V]{
		m:    m,
		keys: make([]K, m),
		vals: make([]V, m),
	}
	return ht
}

// Put inserts a key-value pair into the hash table
func (ht *LinearProbingHashTable[K, V]) Put(key K, value V) {
	// Resize if more than half full
	if ht.n >= ht.m/2 {
		ht.resize(2 * ht.m)
	}

	i := ht.hash(key)
	var zero K
	for ht.keys[i] != zero {
		if ht.keys[i] == key {
			ht.vals[i] = value
			return
		}
		i = (i + 1) % ht.m
	}
	ht.keys[i] = key
	ht.vals[i] = value
	ht.n++
}

// Get retrieves the value associated with the given key
func (ht *LinearProbingHashTable[K, V]) Get(key K) (V, bool) {
	var zero K
	for i := ht.hash(key); ht.keys[i] != zero; i = (i + 1) % ht.m {
		if ht.keys[i] == key {
			return ht.vals[i], true
		}
	}
	var zeroV V
	return zeroV, false
}

// Delete removes a key and its associated value from the hash table
func (ht *LinearProbingHashTable[K, V]) Delete(key K) {
	var zero K
	if _, ok := ht.Get(key); !ok {
		return
	}

	// Find position of key
	i := ht.hash(key)
	for ht.keys[i] != key {
		i = (i + 1) % ht.m
	}

	// Delete key and value
	ht.keys[i] = zero
	var zeroV V
	ht.vals[i] = zeroV

	// Rehash all keys in same cluster
	i = (i + 1) % ht.m
	for ht.keys[i] != zero {
		keyToRehash := ht.keys[i]
		valToRehash := ht.vals[i]
		ht.keys[i] = zero
		ht.vals[i] = zeroV
		ht.n--
		ht.Put(keyToRehash, valToRehash)
		i = (i + 1) % ht.m
	}

	ht.n--

	// Resize if less than 1/8 full
	if ht.n > 0 && ht.n <= ht.m/8 {
		ht.resize(ht.m / 2)
	}
}

// Contains returns true if the hash table contains the given key
func (ht *LinearProbingHashTable[K, V]) Contains(key K) bool {
	_, ok := ht.Get(key)
	return ok
}

// Size returns the number of key-value pairs in the hash table
func (ht *LinearProbingHashTable[K, V]) Size() int {
	return ht.n
}

// IsEmpty returns true if the hash table is empty
func (ht *LinearProbingHashTable[K, V]) IsEmpty() bool {
	return ht.n == 0
}

func (ht *LinearProbingHashTable[K, V]) resize(capacity int) {
	temp := NewLinearProbingHashTable[K, V](capacity)
	var zero K
	for i := 0; i < ht.m; i++ {
		if ht.keys[i] != zero {
			temp.Put(ht.keys[i], ht.vals[i])
		}
	}
	ht.keys = temp.keys
	ht.vals = temp.vals
	ht.m = temp.m
}

func (ht *LinearProbingHashTable[K, V]) hash(key K) int {
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
		h = int(any(key).(int64))
	}
	
	if h < 0 {
		h = -h
	}
	return h % ht.m
}
