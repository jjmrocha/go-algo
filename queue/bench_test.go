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

func BenchmarkPQueueEnqueue(b *testing.B) {
	q := NewPriorityQueue[int](intCmp)
	for b.Loop() {
		q.Enqueue(1)
	}
}

func BenchmarkPQueueEnqueueDequeue(b *testing.B) {
	q := NewPriorityQueue[int](intCmp)
	for b.Loop() {
		q.Enqueue(1)
		q.Dequeue()
	}
}

func BenchmarkSyncPQueueEnqueue(b *testing.B) {
	q := NewSyncPriorityQueue[int](intCmp)
	for b.Loop() {
		q.Enqueue(1)
	}
}

func BenchmarkSyncPQueueEnqueueDequeue(b *testing.B) {
	q := NewSyncPriorityQueue[int](intCmp)
	for b.Loop() {
		q.Enqueue(1)
		q.Dequeue()
	}
}
