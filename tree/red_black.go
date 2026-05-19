// Package tree provides an ordered key-value store backed by a Left-Leaning
// Red-Black BST (Sedgewick, 2008). All operations run in guaranteed O(log n) time.
package tree

import (
	"github.com/jjmrocha/go-algo/sort"
	"github.com/jjmrocha/go-algo/stack"
)

type color bool

const (
	red   color = true
	black color = false
)

// Tree is an ordered symbol table mapping keys of type K to values of type V.
// Keys are ordered by the Comparator supplied to New; equal keys map to a
// single value (last write wins).
type Tree[K, V any] struct {
	root *node[K, V]
	cmp  sort.Comparator[K]
	size int
}

type node[K, V any] struct {
	key   K
	value V
	color color
	left  *node[K, V]
	right *node[K, V]
}

// New returns an empty Tree that orders keys with cmp.
func New[K, V any](cmp sort.Comparator[K]) *Tree[K, V] {
	return &Tree[K, V]{
		cmp: cmp,
	}
}

// Get returns the value associated with key and true, or the zero value
// and false if key is not present.
func (t *Tree[K, V]) Get(key K) (V, bool) {
	var zero V

	if t.size == 0 {
		return zero, false
	}

	node := t.root

	for node != nil {
		switch t.cmp(key, node.key) {
		case sort.Before:
			node = node.left
		case sort.After:
			node = node.right
		default:
			return node.value, true
		}
	}

	return zero, false
}

// Contains reports whether key is present in the tree.
func (t *Tree[K, V]) Contains(key K) bool {
	_, exists := t.Get(key)
	return exists
}

// Len returns the number of key-value pairs in the tree.
func (t *Tree[K, V]) Len() int {
	return t.size
}

// Empty reports whether the tree contains no key-value pairs.
func (t *Tree[K, V]) Empty() bool {
	return t.Len() == 0
}

// Put inserts key with value into the tree. If key already exists its value
// is replaced and Len is unchanged.
func (t *Tree[K, V]) Put(key K, value V) {
	t.root = t.insert(t.root, key, value)
	t.root.color = black
}

func (t *Tree[K, V]) insert(h *node[K, V], key K, value V) *node[K, V] {
	if h == nil {
		t.size++
		return &node[K, V]{
			key:   key,
			value: value,
			color: red,
		}
	}

	switch t.cmp(key, h.key) {
	case sort.Before:
		h.left = t.insert(h.left, key, value)
	case sort.After:
		h.right = t.insert(h.right, key, value)
	default:
		h.value = value
	}

	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}

	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}

	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}

	return h
}

// ToList returns all values in ascending key order.
func (t *Tree[K, V]) ToList() []V {
	result := make([]V, 0, t.size)
	s := stack.New[*node[K, V]]()
	node := t.root

	for node != nil || !s.Empty() {
		for node != nil {
			s.Push(node)
			node = node.left
		}
		node, _ = s.Pop()
		result = append(result, node.value)
		node = node.right
	}

	return result
}

// Min returns the minimum key in the tree and true, or the zero value and false if the tree is empty.
func (t *Tree[K, V]) Min() (K, bool) {
	var zero K

	if t.size == 0 {
		return zero, false
	}

	node := t.root

	for node.left != nil {
		node = node.left
	}

	return node.key, true
}

// Max returns the maximum key in the tree and true, or the zero value and false if the tree is empty.
func (t *Tree[K, V]) Max() (K, bool) {
	var zero K

	if t.size == 0 {
		return zero, false
	}

	node := t.root

	for node.right != nil {
		node = node.right
	}

	return node.key, true
}

// Rank returns the number of keys in the tree that are less than key.
func (t *Tree[K, V]) Rank(key K) int {
	return t.rankOf(t.root, key)
}

func (t *Tree[K, V]) rankOf(node *node[K, V], key K) int {
	if node == nil {
		return 0
	}

	switch t.cmp(key, node.key) {
	case sort.Before:
		return t.rankOf(node.left, key)
	case sort.After:
		return 1 + sizeOf(node.left) + t.rankOf(node.right, key)
	default:
		return sizeOf(node.left)
	}
}

// Select returns the key of rank r (the key such that exactly rank keys in the tree are less than it) and true,
// or the zero value and false if r is out of bounds.
func (t *Tree[K, V]) Select(r int) (K, bool) {
	var zero K

	if r < 0 || r >= t.size {
		return zero, false
	}

	node := selectNode(t.root, r)
	if node == nil {
		return zero, false
	}

	return node.key, true
}

func selectNode[K, V any](node *node[K, V], rank int) *node[K, V] {
	if node == nil {
		return nil
	}

	leftSize := sizeOf(node.left)

	switch {
	case rank < leftSize:
		return selectNode(node.left, rank)
	case rank > leftSize:
		return selectNode(node.right, rank-leftSize-1)
	default:
		return node
	}
}

func sizeOf[K, V any](node *node[K, V]) int {
	if node == nil {
		return 0
	}

	return 1 + sizeOf(node.left) + sizeOf(node.right)
}

func isRed[K, V any](node *node[K, V]) bool {
	if node == nil {
		return false
	}

	return node.color == red
}

func rotateLeft[K, V any](h *node[K, V]) *node[K, V] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = red
	return x
}

func rotateRight[K, V any](h *node[K, V]) *node[K, V] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = red
	return x
}

func flipColors[K, V any](h *node[K, V]) {
	h.color = red
	h.left.color = black
	h.right.color = black
}
