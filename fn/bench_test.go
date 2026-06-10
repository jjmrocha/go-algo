package fn

import (
	"slices"
	"testing"
)

const (
	benchEven = "even"
	benchOdd  = "odd"
)

var n1000 = makeInts(1000)

func makeInts(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func isEven(v int) bool { return v%2 == 0 }
func double(v int) int  { return v * 2 }
func add(a, b int) int  { return a + b }
func parity(v int) string {
	if v%2 == 0 {
		return benchEven
	}
	return benchOdd
}

func BenchmarkFilter(b *testing.B) {
	for b.Loop() {
		_ = Filter(n1000, isEven)
	}
}

func BenchmarkFilterSeq(b *testing.B) {
	for b.Loop() {
		_ = slices.Collect(FilterSeq(slices.Values(n1000), isEven))
	}
}

func BenchmarkMap(b *testing.B) {
	for b.Loop() {
		_ = Map(n1000, double)
	}
}

func BenchmarkMapSeq(b *testing.B) {
	for b.Loop() {
		_ = slices.Collect(MapSeq(slices.Values(n1000), double))
	}
}

func BenchmarkFold(b *testing.B) {
	for b.Loop() {
		_ = Fold(n1000, 0, add)
	}
}

func withDuplicates() []int {
	input := make([]int, 0, len(n1000)+500)
	input = append(input, n1000...)
	input = append(input, n1000[:500]...)
	return input
}

func BenchmarkDistinct(b *testing.B) {
	input := withDuplicates()
	for b.Loop() {
		_ = Distinct(input)
	}
}

func BenchmarkDistinctSeq(b *testing.B) {
	input := withDuplicates()
	for b.Loop() {
		_ = slices.Collect(DistinctSeq(slices.Values(input)))
	}
}

func BenchmarkGroupBy(b *testing.B) {
	for b.Loop() {
		_ = GroupBy(n1000, parity)
	}
}

func BenchmarkFind(b *testing.B) {
	for b.Loop() {
		_, _ = Find(n1000, func(v int) bool { return v == 999 }) // worst-case: last element
	}
}
