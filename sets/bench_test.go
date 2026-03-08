package sets

import "testing"

var sink bool

func BenchmarkSetAdd(b *testing.B) {
	s := New[int]()
	for b.Loop() {
		s.Add(1)
	}
}

func BenchmarkSetContainsHit(b *testing.B) {
	s := Of([]int{1, 2, 3, 4, 5})
	for b.Loop() {
		sink = s.Contains(3)
	}
}

func BenchmarkSetContainsMiss(b *testing.B) {
	s := Of([]int{1, 2, 3, 4, 5})
	for b.Loop() {
		sink = s.Contains(99)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	a := Of(makeRange(0, 500))
	c := Of(makeRange(250, 750))
	for b.Loop() {
		_ = a.Union(c)
	}
}

func BenchmarkSetIntersection(b *testing.B) {
	a := Of(makeRange(0, 500))
	c := Of(makeRange(250, 750))
	for b.Loop() {
		_ = a.Intersection(c)
	}
}

func BenchmarkSetDifference(b *testing.B) {
	a := Of(makeRange(0, 500))
	c := Of(makeRange(250, 750))
	for b.Loop() {
		_ = a.Difference(c)
	}
}

func makeRange(lo, hi int) []int {
	out := make([]int, hi-lo)
	for i := range out {
		out[i] = lo + i
	}
	return out
}
