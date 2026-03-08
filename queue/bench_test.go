package queue

import "testing"

func BenchmarkQueueEnqueue(b *testing.B) {
	q := New[int]()
	for b.Loop() {
		q.Enqueue(1)
	}
}

func BenchmarkQueueEnqueueDequeue(b *testing.B) {
	q := New[int]()
	for b.Loop() {
		q.Enqueue(1)
		q.Dequeue()
	}
}

func BenchmarkSyncQueueEnqueue(b *testing.B) {
	q := NewSyncQueue[int]()
	for b.Loop() {
		q.Enqueue(1)
	}
}

func BenchmarkSyncQueueEnqueueDequeue(b *testing.B) {
	q := NewSyncQueue[int]()
	for b.Loop() {
		q.Enqueue(1)
		q.Dequeue()
	}
}

func BenchmarkBlockingQueueEnqueueDequeue(b *testing.B) {
	q, _ := NewBlockingQueue[int](1)
	for b.Loop() {
		q.Enqueue(1)
		q.Dequeue()
	}
}
