package quickunion

import "testing"

func BenchmarkQuickUnionFind(b *testing.B) {
	u := New(1000)
	for b.Loop() {
		_, _ = u.Find(500)
	}
}

func BenchmarkQuickUnionUnion(b *testing.B) {
	for b.Loop() {
		u := New(1000)
		for i := range 999 {
			u.Union(i, i+1) //nolint:errcheck
		}
	}
}

func BenchmarkQuickUnionConnected(b *testing.B) {
	u := New(1000)
	for i := range 999 {
		u.Union(i, i+1) //nolint:errcheck
	}
	b.ResetTimer()
	for b.Loop() {
		u.Connected(0, 999) //nolint:errcheck
	}
}
