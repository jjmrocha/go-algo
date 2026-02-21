package trees

import "golang.org/x/exp/constraints"

const (
	RED   = true
	BLACK = false
)

// RedBlackBST is a left-leaning red-black tree implementation
type RedBlackBST[K constraints.Ordered, V any] struct {
	root *rbNode[K, V]
}

type rbNode[K constraints.Ordered, V any] struct {
	key   K
	value V
	left  *rbNode[K, V]
	right *rbNode[K, V]
	color bool
	size  int
}

// NewRedBlackBST creates a new red-black binary search tree
func NewRedBlackBST[K constraints.Ordered, V any]() *RedBlackBST[K, V] {
	return &RedBlackBST[K, V]{}
}

// Put inserts a key-value pair into the tree
func (rbt *RedBlackBST[K, V]) Put(key K, value V) {
	rbt.root = rbt.put(rbt.root, key, value)
	rbt.root.color = BLACK
}

func (rbt *RedBlackBST[K, V]) put(h *rbNode[K, V], key K, value V) *rbNode[K, V] {
	if h == nil {
		return &rbNode[K, V]{key: key, value: value, color: RED, size: 1}
	}

	if key < h.key {
		h.left = rbt.put(h.left, key, value)
	} else if key > h.key {
		h.right = rbt.put(h.right, key, value)
	} else {
		h.value = value
	}

	// Fix right-leaning reds
	if rbt.isRed(h.right) && !rbt.isRed(h.left) {
		h = rbt.rotateLeft(h)
	}
	// Fix two reds in a row
	if rbt.isRed(h.left) && rbt.isRed(h.left.left) {
		h = rbt.rotateRight(h)
	}
	// Split 4-nodes
	if rbt.isRed(h.left) && rbt.isRed(h.right) {
		rbt.flipColors(h)
	}

	h.size = 1 + rbt.size(h.left) + rbt.size(h.right)
	return h
}

// Get retrieves the value associated with the given key
func (rbt *RedBlackBST[K, V]) Get(key K) (V, bool) {
	n := rbt.get(rbt.root, key)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.value, true
}

func (rbt *RedBlackBST[K, V]) get(n *rbNode[K, V], key K) *rbNode[K, V] {
	if n == nil {
		return nil
	}
	if key < n.key {
		return rbt.get(n.left, key)
	} else if key > n.key {
		return rbt.get(n.right, key)
	}
	return n
}

// Contains returns true if the tree contains the given key
func (rbt *RedBlackBST[K, V]) Contains(key K) bool {
	_, ok := rbt.Get(key)
	return ok
}

// IsEmpty returns true if the tree is empty
func (rbt *RedBlackBST[K, V]) IsEmpty() bool {
	return rbt.root == nil
}

// Size returns the number of key-value pairs in the tree
func (rbt *RedBlackBST[K, V]) Size() int {
	return rbt.size(rbt.root)
}

func (rbt *RedBlackBST[K, V]) size(n *rbNode[K, V]) int {
	if n == nil {
		return 0
	}
	return n.size
}

func (rbt *RedBlackBST[K, V]) isRed(n *rbNode[K, V]) bool {
	if n == nil {
		return false
	}
	return n.color == RED
}

func (rbt *RedBlackBST[K, V]) rotateLeft(h *rbNode[K, V]) *rbNode[K, V] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = 1 + rbt.size(h.left) + rbt.size(h.right)
	return x
}

func (rbt *RedBlackBST[K, V]) rotateRight(h *rbNode[K, V]) *rbNode[K, V] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = 1 + rbt.size(h.left) + rbt.size(h.right)
	return x
}

func (rbt *RedBlackBST[K, V]) flipColors(h *rbNode[K, V]) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}
