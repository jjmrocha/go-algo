// Package quickunion implements a weighted quick-union data structure with
// path compression for efficiently tracking disjoint sets.
package quickunion

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrIndexOutOfRange is returned by Find when the given index is outside
// the valid range [0, size).
var ErrIndexOutOfRange = errors.New("index out of bounds")

// QuickUnion is a disjoint-set data structure that tracks a
// collection of elements partitioned into non-overlapping sets. It uses
// union by weight and half-path compression to achieve near-constant
// time for each operation.
type QuickUnion struct {
	sets   int
	parent []int
	weight []int
}

// New creates a QuickUnion with size elements, each initially in its own set.
func New(size int) *QuickUnion {
	p := make([]int, size)
	w := make([]int, size)

	for i := range size {
		p[i] = i
		w[i] = 1
	}

	return &QuickUnion{
		sets:   size,
		parent: p,
		weight: w,
	}
}

// Len returns the number of disjoint sets.
func (u *QuickUnion) Len() int {
	return u.sets
}

// Find returns the root of the set containing p, applying half-path
// compression along the way. It returns ErrIndexOutOfRange if p is out of
// bounds.
func (u *QuickUnion) Find(p int) (int, error) {
	if p < 0 || p >= len(u.parent) {
		return 0, ErrIndexOutOfRange
	}

	for p != u.parent[p] {
		parent := u.parent[p]
		u.parent[p] = u.parent[parent]
		p = u.parent[p]
	}

	return u.parent[p], nil
}

// Union merges the sets containing p and q using union by weight.
// It returns nil if the sets were successfully merged or were already
// connected, and ErrIndexOutOfRange if either index is out of bounds.
func (u *QuickUnion) Union(p, q int) error {
	rp, err := u.Find(p)
	if err != nil {
		return err
	}

	rq, err := u.Find(q)
	if err != nil {
		return err
	}

	if rp == rq {
		return nil
	}

	wp := u.weight[rp]
	wq := u.weight[rq]

	if wp < wq {
		u.parent[rp] = rq
		u.weight[rq] += wp
	} else {
		u.parent[rq] = rp
		u.weight[rp] += wq
	}

	u.sets--
	return nil
}

// Connected reports whether p and q belong to the same set.
// It returns ErrIndexOutOfRange if either index is out of bounds.
func (u *QuickUnion) Connected(p, q int) (bool, error) {
	rp, err := u.Find(p)
	if err != nil {
		return false, err
	}

	rq, err := u.Find(q)
	if err != nil {
		return false, err
	}

	return rp == rq, nil
}

// String returns a human-readable representation of the parent links.
// Each non-root element is printed as "parent <- child".
func (u *QuickUnion) String() string {
	var connections []string

	for i, p := range u.parent {
		if i != p {
			s := fmt.Sprintf("%d <- %d", p, i)
			connections = append(connections, s)
		}
	}

	slices.Sort(connections)

	return strings.Join(connections, "\n")
}
