package trees

import "golang.org/x/exp/constraints"

// BTree is a B-tree implementation
type BTree[K constraints.Ordered, V any] struct {
	root *bNode[K, V]
	t    int // minimum degree (minimum number of keys is t-1)
}

type bNode[K constraints.Ordered, V any] struct {
	keys     []K
	values   []V
	children []*bNode[K, V]
	leaf     bool
}

// NewBTree creates a new B-tree with the given minimum degree
// t >= 2, each node can have at most 2t-1 keys and 2t children
func NewBTree[K constraints.Ordered, V any](t int) *BTree[K, V] {
	if t < 2 {
		t = 2
	}
	return &BTree[K, V]{
		root: &bNode[K, V]{
			keys:   make([]K, 0),
			values: make([]V, 0),
			leaf:   true,
		},
		t: t,
	}
}

// Search searches for a key in the B-tree
func (bt *BTree[K, V]) Search(key K) (V, bool) {
	return bt.search(bt.root, key)
}

func (bt *BTree[K, V]) search(n *bNode[K, V], key K) (V, bool) {
	i := 0
	for i < len(n.keys) && key > n.keys[i] {
		i++
	}

	if i < len(n.keys) && key == n.keys[i] {
		return n.values[i], true
	}

	if n.leaf {
		var zero V
		return zero, false
	}

	return bt.search(n.children[i], key)
}

// Insert inserts a key-value pair into the B-tree
func (bt *BTree[K, V]) Insert(key K, value V) {
	r := bt.root
	if len(r.keys) == 2*bt.t-1 {
		s := &bNode[K, V]{
			keys:     make([]K, 0),
			values:   make([]V, 0),
			children: make([]*bNode[K, V], 0),
			leaf:     false,
		}
		bt.root = s
		s.children = append(s.children, r)
		bt.splitChild(s, 0)
		bt.insertNonFull(s, key, value)
	} else {
		bt.insertNonFull(r, key, value)
	}
}

func (bt *BTree[K, V]) insertNonFull(n *bNode[K, V], key K, value V) {
	i := len(n.keys) - 1

	if n.leaf {
		n.keys = append(n.keys, key)
		n.values = append(n.values, value)
		for i >= 0 && key < n.keys[i] {
			n.keys[i+1] = n.keys[i]
			n.values[i+1] = n.values[i]
			i--
		}
		n.keys[i+1] = key
		n.values[i+1] = value
	} else {
		for i >= 0 && key < n.keys[i] {
			i--
		}
		i++
		if len(n.children[i].keys) == 2*bt.t-1 {
			bt.splitChild(n, i)
			if key > n.keys[i] {
				i++
			}
		}
		bt.insertNonFull(n.children[i], key, value)
	}
}

func (bt *BTree[K, V]) splitChild(parent *bNode[K, V], i int) {
	t := bt.t
	full := parent.children[i]
	newNode := &bNode[K, V]{
		keys:   make([]K, 0),
		values: make([]V, 0),
		leaf:   full.leaf,
	}

	// Move the second half of keys and values to new node
	newNode.keys = append(newNode.keys, full.keys[t:]...)
	newNode.values = append(newNode.values, full.values[t:]...)
	full.keys = full.keys[:t-1]
	full.values = full.values[:t-1]

	// If not a leaf, move children too
	if !full.leaf {
		newNode.children = make([]*bNode[K, V], 0)
		newNode.children = append(newNode.children, full.children[t:]...)
		full.children = full.children[:t]
	}

	// Insert middle key into parent
	parent.keys = append(parent.keys, full.keys[t-1])
	parent.values = append(parent.values, full.values[t-1])
	for j := len(parent.keys) - 1; j > i; j-- {
		parent.keys[j] = parent.keys[j-1]
		parent.values[j] = parent.values[j-1]
	}
	parent.keys[i] = full.keys[t-1]
	parent.values[i] = full.values[t-1]

	// Insert new child into parent
	parent.children = append(parent.children, newNode)
	for j := len(parent.children) - 1; j > i+1; j-- {
		parent.children[j] = parent.children[j-1]
	}
	parent.children[i+1] = newNode
}

// IsEmpty returns true if the B-tree is empty
func (bt *BTree[K, V]) IsEmpty() bool {
	return len(bt.root.keys) == 0
}
