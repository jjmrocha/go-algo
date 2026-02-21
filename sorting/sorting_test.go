package sorting

import (
	"reflect"
	"testing"
)

func TestSelectionSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	expected := []int{11, 12, 22, 25, 64}
	SelectionSort(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, got %v", expected, arr)
	}
}

func TestInsertionSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	expected := []int{11, 12, 22, 25, 64}
	InsertionSort(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, got %v", expected, arr)
	}
}

func TestMergeSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	expected := []int{11, 12, 22, 25, 64}
	MergeSort(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, got %v", expected, arr)
	}
}

func TestQuickSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	expected := []int{11, 12, 22, 25, 64}
	QuickSort(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, got %v", expected, arr)
	}
}

func TestHeapSort(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	expected := []int{11, 12, 22, 25, 64}
	HeapSort(arr)
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("Expected %v, got %v", expected, arr)
	}
}

func TestQuickSelect(t *testing.T) {
	arr := []int{64, 25, 12, 22, 11}
	// Find 3rd smallest element (0-indexed, so k=2)
	result := QuickSelect(arr, 2)
	if result != 22 {
		t.Errorf("Expected 22, got %d", result)
	}
}
