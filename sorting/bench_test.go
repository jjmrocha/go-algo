package sorting

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

var benchSizes = []int{1_000, 10_000, 100_000, 1_000_000}
var benchDists = []string{"random", "sorted", "reverse", "dupes15"}

func intAsc(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func makeInts(n int, dist string) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	switch dist {
	case "random":
		rand.Shuffle(n, func(i, j int) { s[i], s[j] = s[j], s[i] })
	case "reverse":
		for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
	case "dupes15":
		unique := int(float64(n) * 0.85)
		for i := range s {
			s[i] = i % unique
		}
		rand.Shuffle(n, func(i, j int) { s[i], s[j] = s[j], s[i] })
	}
	return s
}

func BenchmarkMerge(b *testing.B) {
	for _, n := range benchSizes {
		for _, dist := range benchDists {
			b.Run(fmt.Sprintf("n=%d/%s", n, dist), func(b *testing.B) {
				base := makeInts(n, dist)
				for b.Loop() {
					input := make([]int, n)
					copy(input, base)
					Merge(input, intAsc)
				}
			})
		}
	}
}

func BenchmarkQuick(b *testing.B) {
	for _, n := range benchSizes {
		for _, dist := range benchDists {
			b.Run(fmt.Sprintf("n=%d/%s", n, dist), func(b *testing.B) {
				base := makeInts(n, dist)
				for b.Loop() {
					input := make([]int, n)
					copy(input, base)
					Quick(input, intAsc)
				}
			})
		}
	}
}
