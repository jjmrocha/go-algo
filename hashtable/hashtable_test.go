package hashtable

import "testing"

func TestSeparateChainingHashTable(t *testing.T) {
	ht := NewSeparateChainingHashTable[int, string](5)
	
	if !ht.IsEmpty() {
		t.Error("New hash table should be empty")
	}
	
	ht.Put(1, "one")
	ht.Put(2, "two")
	ht.Put(3, "three")
	
	if ht.Size() != 3 {
		t.Errorf("Expected size 3, got %d", ht.Size())
	}
	
	val, ok := ht.Get(2)
	if !ok || val != "two" {
		t.Errorf("Expected 'two', got '%s'", val)
	}
	
	if !ht.Contains(3) {
		t.Error("Hash table should contain key 3")
	}
	
	ht.Delete(2)
	if ht.Contains(2) {
		t.Error("Hash table should not contain key 2 after deletion")
	}
}

func TestLinearProbingHashTable(t *testing.T) {
	ht := NewLinearProbingHashTable[int, string](16)
	
	if !ht.IsEmpty() {
		t.Error("New hash table should be empty")
	}
	
	ht.Put(1, "one")
	ht.Put(2, "two")
	ht.Put(3, "three")
	
	if ht.Size() != 3 {
		t.Errorf("Expected size 3, got %d", ht.Size())
	}
	
	val, ok := ht.Get(2)
	if !ok || val != "two" {
		t.Errorf("Expected 'two', got '%s'", val)
	}
	
	if !ht.Contains(3) {
		t.Error("Hash table should contain key 3")
	}
	
	ht.Delete(2)
	if ht.Contains(2) {
		t.Error("Hash table should not contain key 2 after deletion")
	}
}
