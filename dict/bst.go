package dict

import "golang.org/x/exp/constraints"

// BST is a binary search tree implementation of a symbol table
type BST[K constraints.Ordered, V any] struct {
	root *node[K, V]
}

type node[K constraints.Ordered, V any] struct {
	key   K
	value V
	left  *node[K, V]
	right *node[K, V]
	size  int
}

// NewBST creates a new binary search tree
func NewBST[K constraints.Ordered, V any]() *BST[K, V] {
	return &BST[K, V]{}
}

// Put inserts a key-value pair into the tree
func (bst *BST[K, V]) Put(key K, value V) {
	bst.root = bst.put(bst.root, key, value)
}

func (bst *BST[K, V]) put(n *node[K, V], key K, value V) *node[K, V] {
	if n == nil {
		return &node[K, V]{key: key, value: value, size: 1}
	}
	if key < n.key {
		n.left = bst.put(n.left, key, value)
	} else if key > n.key {
		n.right = bst.put(n.right, key, value)
	} else {
		n.value = value
	}
	n.size = 1 + bst.size(n.left) + bst.size(n.right)
	return n
}

// Get retrieves the value associated with the given key
func (bst *BST[K, V]) Get(key K) (V, bool) {
	n := bst.get(bst.root, key)
	if n == nil {
		var zero V
		return zero, false
	}
	return n.value, true
}

func (bst *BST[K, V]) get(n *node[K, V], key K) *node[K, V] {
	if n == nil {
		return nil
	}
	if key < n.key {
		return bst.get(n.left, key)
	} else if key > n.key {
		return bst.get(n.right, key)
	}
	return n
}

// Delete removes a key and its associated value from the tree
func (bst *BST[K, V]) Delete(key K) {
	bst.root = bst.delete(bst.root, key)
}

func (bst *BST[K, V]) delete(n *node[K, V], key K) *node[K, V] {
	if n == nil {
		return nil
	}
	if key < n.key {
		n.left = bst.delete(n.left, key)
	} else if key > n.key {
		n.right = bst.delete(n.right, key)
	} else {
		if n.right == nil {
			return n.left
		}
		if n.left == nil {
			return n.right
		}
		t := n
		n = bst.min(t.right)
		n.right = bst.deleteMin(t.right)
		n.left = t.left
	}
	n.size = 1 + bst.size(n.left) + bst.size(n.right)
	return n
}

func (bst *BST[K, V]) deleteMin(n *node[K, V]) *node[K, V] {
	if n.left == nil {
		return n.right
	}
	n.left = bst.deleteMin(n.left)
	n.size = 1 + bst.size(n.left) + bst.size(n.right)
	return n
}

func (bst *BST[K, V]) min(n *node[K, V]) *node[K, V] {
	if n.left == nil {
		return n
	}
	return bst.min(n.left)
}

// Contains returns true if the tree contains the given key
func (bst *BST[K, V]) Contains(key K) bool {
	_, ok := bst.Get(key)
	return ok
}

// IsEmpty returns true if the tree is empty
func (bst *BST[K, V]) IsEmpty() bool {
	return bst.Size() == 0
}

// Size returns the number of key-value pairs in the tree
func (bst *BST[K, V]) Size() int {
	return bst.size(bst.root)
}

func (bst *BST[K, V]) size(n *node[K, V]) int {
	if n == nil {
		return 0
	}
	return n.size
}
