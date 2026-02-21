package trees

import (
	"math"
)

// Point2D represents a 2D point
type Point2D struct {
	X, Y float64
}

// KDTree is a 2D tree for efficient 2D point search
type KDTree struct {
	root *kdNode
	size int
}

type kdNode struct {
	point Point2D
	left  *kdNode
	right *kdNode
}

// NewKDTree creates a new kd-tree
func NewKDTree() *KDTree {
	return &KDTree{}
}

// Insert adds a point to the kd-tree
func (kd *KDTree) Insert(p Point2D) {
	kd.root = kd.insert(kd.root, p, true)
}

func (kd *KDTree) insert(n *kdNode, p Point2D, useX bool) *kdNode {
	if n == nil {
		kd.size++
		return &kdNode{point: p}
	}

	if p.X == n.point.X && p.Y == n.point.Y {
		return n
	}

	if useX {
		if p.X < n.point.X {
			n.left = kd.insert(n.left, p, !useX)
		} else {
			n.right = kd.insert(n.right, p, !useX)
		}
	} else {
		if p.Y < n.point.Y {
			n.left = kd.insert(n.left, p, !useX)
		} else {
			n.right = kd.insert(n.right, p, !useX)
		}
	}

	return n
}

// Contains returns true if the tree contains the given point
func (kd *KDTree) Contains(p Point2D) bool {
	return kd.contains(kd.root, p, true)
}

func (kd *KDTree) contains(n *kdNode, p Point2D, useX bool) bool {
	if n == nil {
		return false
	}

	if p.X == n.point.X && p.Y == n.point.Y {
		return true
	}

	if useX {
		if p.X < n.point.X {
			return kd.contains(n.left, p, !useX)
		}
		return kd.contains(n.right, p, !useX)
	}

	if p.Y < n.point.Y {
		return kd.contains(n.left, p, !useX)
	}
	return kd.contains(n.right, p, !useX)
}

// Nearest finds the nearest point to the given point
func (kd *KDTree) Nearest(p Point2D) *Point2D {
	if kd.root == nil {
		return nil
	}
	nearest := kd.root.point
	minDist := kd.distance(p, nearest)
	return kd.nearest(kd.root, p, &nearest, &minDist, true)
}

func (kd *KDTree) nearest(n *kdNode, p Point2D, best *Point2D, bestDist *float64, useX bool) *Point2D {
	if n == nil {
		return best
	}

	dist := kd.distance(p, n.point)
	if dist < *bestDist {
		*best = n.point
		*bestDist = dist
	}

	var first, second *kdNode
	var axisDist float64

	if useX {
		axisDist = p.X - n.point.X
		if axisDist < 0 {
			first = n.left
			second = n.right
		} else {
			first = n.right
			second = n.left
		}
	} else {
		axisDist = p.Y - n.point.Y
		if axisDist < 0 {
			first = n.left
			second = n.right
		} else {
			first = n.right
			second = n.left
		}
	}

	best = kd.nearest(first, p, best, bestDist, !useX)

	// Check if we need to search the other side
	if axisDist*axisDist < *bestDist {
		best = kd.nearest(second, p, best, bestDist, !useX)
	}

	return best
}

// RangeSearch finds all points within a rectangle
func (kd *KDTree) RangeSearch(xMin, yMin, xMax, yMax float64) []Point2D {
	points := make([]Point2D, 0)
	kd.rangeSearch(kd.root, xMin, yMin, xMax, yMax, &points, true)
	return points
}

func (kd *KDTree) rangeSearch(n *kdNode, xMin, yMin, xMax, yMax float64, points *[]Point2D, useX bool) {
	if n == nil {
		return
	}

	if n.point.X >= xMin && n.point.X <= xMax && n.point.Y >= yMin && n.point.Y <= yMax {
		*points = append(*points, n.point)
	}

	if useX {
		if xMin <= n.point.X {
			kd.rangeSearch(n.left, xMin, yMin, xMax, yMax, points, !useX)
		}
		if xMax >= n.point.X {
			kd.rangeSearch(n.right, xMin, yMin, xMax, yMax, points, !useX)
		}
	} else {
		if yMin <= n.point.Y {
			kd.rangeSearch(n.left, xMin, yMin, xMax, yMax, points, !useX)
		}
		if yMax >= n.point.Y {
			kd.rangeSearch(n.right, xMin, yMin, xMax, yMax, points, !useX)
		}
	}
}

// Size returns the number of points in the tree
func (kd *KDTree) Size() int {
	return kd.size
}

// IsEmpty returns true if the tree is empty
func (kd *KDTree) IsEmpty() bool {
	return kd.size == 0
}

func (kd *KDTree) distance(p1, p2 Point2D) float64 {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}
