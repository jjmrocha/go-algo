// Package union implements a weighted quick-union data structure with
// path compression for efficiently tracking disjoint sets.
package union

import (
	"fmt"
	"slices"
	"strings"
)

// IndexOutOfRange is returned by Find when the given index is outside
// the valid range [0, size).
const IndexOutOfRange = -1

// QuickUnion is a disjoint-set data structure that tracks a
// collection of elements partitioned into non-overlapping sets. It uses
// union by weight and half-path compression to achieve near-constant
// time for each operation.
type QuickUnion struct {
	sets    int
	parents []int
	weights []int
}

// New creates a QuickUnion with size elements, each initially in its own set.
func New(size int) *QuickUnion {
	p := make([]int, size)
	w := make([]int, size)

	for i := 0; i < size; i++ {
		p[i] = i
		w[i] = 1
	}

	return &QuickUnion{
		sets:    size,
		parents: p,
		weights: w,
	}
}

// Count returns the number of disjoint sets.
func (u *QuickUnion) Count() int {
	return u.sets
}

// Find returns the root of the set containing p, applying half-path
// compression along the way. It returns IndexOutOfRange if p is out of
// bounds.
func (u *QuickUnion) Find(p int) int {
	if p < 0 || p >= len(u.parents) {
		return IndexOutOfRange
	}

	for p != u.parents[p] {
		parent := u.parents[p]
		u.parents[p] = u.parents[parent]
		p = parent
	}

	return u.parents[p]
}

// Union merges the sets containing p and q using union by weight.
// It returns true if the two elements were in different sets and were
// successfully merged, and false if they were already connected or either
// index is out of bounds.
func (u *QuickUnion) Union(p, q int) bool {
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
		u.parents[rp] = rq
		u.weights[rq] += wp
	} else {
		u.parents[rq] = rp
		u.weights[rp] += wq
	}

	u.sets--
	return true
}

// IsConnected reports whether p and q belong to the same set.
// It returns false if either index is out of bounds.
func (u *QuickUnion) IsConnected(p, q int) bool {
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

// String returns a human-readable representation of the parents links.
// Each non-root element is printed as "parents <- child".
func (u *QuickUnion) String() string {
	connections := make([]string, 0)

	for i, p := range u.parents {
		if i != p {
			s := fmt.Sprintf("%d <- %d", p, i)
			connections = append(connections, s)
		}
	}

	slices.Sort(connections)

	return strings.Join(connections, "\n")
}
