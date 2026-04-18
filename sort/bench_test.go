package sort

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

var benchSizes = []int{1_000, 10_000, 100_000}
var benchDists = []string{"random", "sorted", "reverse"}

func intAsc(a, b int) int {
	if a < b {
		return Before
	}
	if a > b {
		return After
	}
	return Equal
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
	}
	return s
}

func BenchmarkMerge(b *testing.B) {
	for _, n := range benchSizes {
		for _, dist := range benchDists {
			b.Run(fmt.Sprintf("n=%d/%s", n, dist), func(b *testing.B) {
				input := makeInts(n, dist)
				for b.Loop() {
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
