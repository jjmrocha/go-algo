package token

import (
	"testing"
	"uuid"
)

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = New()
	}
}

func BenchmarkFromUUID(b *testing.B) {
	u := uuid.MustParse("41c9ad60-0cab-4e1b-afc8-cf97fbb94662")
	b.ReportAllocs()
	for b.Loop() {
		_ = FromUUID(u)
	}
}

func BenchmarkValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = Valid("3w7nni025418b96lydzqxyb8i")
	}
}
