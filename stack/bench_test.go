package stack

import "testing"

func BenchmarkStackPush(b *testing.B) {
	s := New[int]()
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkStackPushPop(b *testing.B) {
	s := New[int]()
	for b.Loop() {
		s.Push(1)
		s.Pop()
	}
}

func BenchmarkStackPeek(b *testing.B) {
	s := New[int]()
	s.Push(42)
	for b.Loop() {
		s.Peek()
	}
}

func BenchmarkSyncStackPush(b *testing.B) {
	s := NewSyncStack[int]()
	for b.Loop() {
		s.Push(1)
	}
}

func BenchmarkSyncStackPushPop(b *testing.B) {
	s := NewSyncStack[int]()
	for b.Loop() {
		s.Push(1)
		s.Pop()
	}
}
