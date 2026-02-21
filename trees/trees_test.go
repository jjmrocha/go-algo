package trees

import "testing"

func TestRedBlackBST(t *testing.T) {
	rbt := NewRedBlackBST[int, string]()
	
	if !rbt.IsEmpty() {
		t.Error("New Red-Black BST should be empty")
	}
	
	rbt.Put(5, "five")
	rbt.Put(3, "three")
	rbt.Put(7, "seven")
	rbt.Put(1, "one")
	
	if rbt.Size() != 4 {
		t.Errorf("Expected size 4, got %d", rbt.Size())
	}
	
	val, ok := rbt.Get(3)
	if !ok || val != "three" {
		t.Errorf("Expected 'three', got '%s'", val)
	}
	
	if !rbt.Contains(7) {
		t.Error("Red-Black BST should contain key 7")
	}
}

func TestBTree(t *testing.T) {
	bt := NewBTree[int, string](3)
	
	if !bt.IsEmpty() {
		t.Error("New B-tree should be empty")
	}
	
	bt.Insert(5, "five")
	bt.Insert(3, "three")
	bt.Insert(7, "seven")
	bt.Insert(1, "one")
	
	val, ok := bt.Search(3)
	if !ok || val != "three" {
		t.Errorf("Expected 'three', got '%s'", val)
	}
	
	val, ok = bt.Search(7)
	if !ok || val != "seven" {
		t.Errorf("Expected 'seven', got '%s'", val)
	}
}

func TestKDTree(t *testing.T) {
	kd := NewKDTree()
	
	if !kd.IsEmpty() {
		t.Error("New KD-tree should be empty")
	}
	
	p1 := Point2D{X: 1.0, Y: 2.0}
	p2 := Point2D{X: 3.0, Y: 4.0}
	p3 := Point2D{X: 5.0, Y: 6.0}
	
	kd.Insert(p1)
	kd.Insert(p2)
	kd.Insert(p3)
	
	if kd.Size() != 3 {
		t.Errorf("Expected size 3, got %d", kd.Size())
	}
	
	if !kd.Contains(p2) {
		t.Error("KD-tree should contain p2")
	}
	
	nearest := kd.Nearest(Point2D{X: 3.1, Y: 4.1})
	if nearest == nil {
		t.Error("Nearest should not be nil")
	} else if nearest.X != p2.X || nearest.Y != p2.Y {
		t.Errorf("Expected nearest to be p2, got %v", nearest)
	}
}
