package unionfind

import "testing"

func TestQuickFind(t *testing.T) {
	qf := NewQuickFind(10)
	
	if qf.Count() != 10 {
		t.Errorf("Expected count 10, got %d", qf.Count())
	}
	
	qf.Union(4, 3)
	qf.Union(3, 8)
	qf.Union(6, 5)
	
	if !qf.Connected(4, 8) {
		t.Error("4 and 8 should be connected")
	}
	
	if qf.Connected(4, 5) {
		t.Error("4 and 5 should not be connected")
	}
	
	if qf.Count() != 7 {
		t.Errorf("Expected count 7, got %d", qf.Count())
	}
}

func TestQuickUnion(t *testing.T) {
	qu := NewQuickUnion(10)
	
	qu.Union(4, 3)
	qu.Union(3, 8)
	qu.Union(6, 5)
	
	if !qu.Connected(4, 8) {
		t.Error("4 and 8 should be connected")
	}
	
	if qu.Connected(4, 5) {
		t.Error("4 and 5 should not be connected")
	}
}

func TestWeightedQuickUnion(t *testing.T) {
	wqu := NewWeightedQuickUnion(10)
	
	wqu.Union(4, 3)
	wqu.Union(3, 8)
	wqu.Union(6, 5)
	
	if !wqu.Connected(4, 8) {
		t.Error("4 and 8 should be connected")
	}
	
	if wqu.Connected(4, 5) {
		t.Error("4 and 5 should not be connected")
	}
}

func TestWeightedQuickUnionPC(t *testing.T) {
	wqupc := NewWeightedQuickUnionPC(10)
	
	wqupc.Union(4, 3)
	wqupc.Union(3, 8)
	wqupc.Union(6, 5)
	
	if !wqupc.Connected(4, 8) {
		t.Error("4 and 8 should be connected")
	}
	
	if wqupc.Connected(4, 5) {
		t.Error("4 and 5 should not be connected")
	}
}
