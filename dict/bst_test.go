package dict

import "testing"

func TestBST(t *testing.T) {
	bst := NewBST[int, string]()
	
	if !bst.IsEmpty() {
		t.Error("New BST should be empty")
	}
	
	bst.Put(5, "five")
	bst.Put(3, "three")
	bst.Put(7, "seven")
	bst.Put(1, "one")
	
	if bst.Size() != 4 {
		t.Errorf("Expected size 4, got %d", bst.Size())
	}
	
	val, ok := bst.Get(3)
	if !ok || val != "three" {
		t.Errorf("Expected 'three', got '%s'", val)
	}
	
	if !bst.Contains(7) {
		t.Error("BST should contain key 7")
	}
	
	bst.Delete(3)
	if bst.Contains(3) {
		t.Error("BST should not contain key 3 after deletion")
	}
}
