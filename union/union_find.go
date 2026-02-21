// Package union implements a weighted quick-union data structure with
// path compression for efficiently tracking disjoint sets.
package union

import (
	"fmt"
)

// IndexOutOfRange is returned by Find when the given index is outside
// the valid range [0, size).
const IndexOutOfRange = -1

// UnionFind is a disjoint-set (union–find) data structure that tracks a
// collection of elements partitioned into non-overlapping sets. It uses
// union by weight and half-path compression to achieve near-constant
// amortised time for each operation.
type UnionFind struct {
	sets    int
	parent  []int
	weights []int
}

// New creates a UnionFind with size elements, each initially in its own set.
func New(size int) *UnionFind {
	p := make([]int, size)
	w := make([]int, size)

	for i := 0; i < size; i++ {
		p[i] = i
		w[i] = 1
	}

	return &UnionFind{
		sets:    size,
		parent:  p,
		weights: w,
	}
}

// Count returns the number of disjoint sets.
func (u *UnionFind) Count() int {
	return u.sets
}

// Find returns the root of the set containing p, applying half-path
// compression along the way. It returns IndexOutOfRange if p is out of
// bounds.
func (u *UnionFind) Find(p int) int {
	if p < 0 || p >= len(u.parent) {
		return IndexOutOfRange
	}

	for p != u.parent[p] {
		parent := u.parent[p]
		u.parent[p] = u.parent[parent]
		p = parent
	}

	return u.parent[p]
}

// Union merges the sets containing p and q using union by weight.
// It returns true if the two elements were in different sets and were
// successfully merged, and false if they were already connected or either
// index is out of bounds.
func (u *UnionFind) Union(p, q int) bool {
	rp := u.Find(p)
	if rp == IndexOutOfRange {
		return false
	}

	rq := u.Find(q)
	if rq == IndexOutOfRange {
		return false
	}

	if rp == rq {
		return false
	}

	wp := u.weights[rp]
	wq := u.weights[rq]

	if wp < wq {
		u.parent[rp] = rq
		u.weights[rq] += wp
	} else {
		u.parent[rq] = rp
		u.weights[rp] += wq
	}

	u.sets--
	return true
}

// IsConnected reports whether p and q belong to the same set.
// It returns false if either index is out of bounds.
func (u *UnionFind) IsConnected(p, q int) bool {
	rp := u.Find(p)
	if rp == IndexOutOfRange {
		return false
	}

	rq := u.Find(q)
	if rq == IndexOutOfRange {
		return false
	}

	return rp == rq
}

// String returns a human-readable representation of the parent links.
// Each non-root element is printed as "parent <- child".
func (u *UnionFind) String() string {
	var s string

	for i, p := range u.parent {
		if i != p {
			s += fmt.Sprintf("%d <- %d \n", p, i)
		}
	}

	return s
}
