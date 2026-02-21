package priorityqueue

import "testing"

func TestMaxPQ(t *testing.T) {
	pq := NewMaxPQ[int]()
	
	if !pq.IsEmpty() {
		t.Error("New priority queue should be empty")
	}
	
	pq.Insert(3)
	pq.Insert(1)
	pq.Insert(5)
	pq.Insert(2)
	
	if pq.Size() != 4 {
		t.Errorf("Expected size 4, got %d", pq.Size())
	}
	
	max, ok := pq.Max()
	if !ok || max != 5 {
		t.Errorf("Expected max to be 5, got %d", max)
	}
	
	val, ok := pq.DelMax()
	if !ok || val != 5 {
		t.Errorf("Expected DelMax to return 5, got %d", val)
	}
	
	val, ok = pq.DelMax()
	if !ok || val != 3 {
		t.Errorf("Expected DelMax to return 3, got %d", val)
	}
}

func TestMinPQ(t *testing.T) {
	pq := NewMinPQ[int]()
	
	pq.Insert(3)
	pq.Insert(1)
	pq.Insert(5)
	pq.Insert(2)
	
	min, ok := pq.Min()
	if !ok || min != 1 {
		t.Errorf("Expected min to be 1, got %d", min)
	}
	
	val, ok := pq.DelMin()
	if !ok || val != 1 {
		t.Errorf("Expected DelMin to return 1, got %d", val)
	}
	
	val, ok = pq.DelMin()
	if !ok || val != 2 {
		t.Errorf("Expected DelMin to return 2, got %d", val)
	}
}
