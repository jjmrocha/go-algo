# go-algo

A comprehensive Go library implementing classic algorithms from Princeton University's Algorithms, Part I course on Coursera.

## Overview

This repository contains Go implementations of fundamental algorithms and data structures, including:
- Union-Find algorithms
- Basic data structures (Stack, Queue)
- Sorting algorithms
- Priority queues
- Associative arrays/dictionaries
- Tree structures (BST, Red-Black BST, B-Trees, kd-Trees)
- Hash tables

## Algorithms and Data Structures

### Union-Find

The `unionfind` package provides four implementations of the union-find data structure:

- **QuickFind**: Simple implementation with O(1) find and O(n) union
- **QuickUnion**: Tree-based implementation with faster union
- **WeightedQuickUnion**: Optimized with weighted trees to keep them balanced
- **WeightedQuickUnionPC**: Further optimized with path compression

```go
import "github.com/jjmrocha/go-algo/unionfind"

uf := unionfind.NewWeightedQuickUnionPC(10)
uf.Union(4, 3)
uf.Union(3, 8)
if uf.Connected(4, 8) {
    // Elements 4 and 8 are connected
}
```

### Stack and Queue

Generic implementations using Go generics:

```go
import (
    "github.com/jjmrocha/go-algo/stack"
    "github.com/jjmrocha/go-algo/queue"
)

// Stack (LIFO)
s := stack.New[int]()
s.Push(1)
s.Push(2)
val, _ := s.Pop() // returns 2

// Queue (FIFO)
q := queue.New[int]()
q.Enqueue(1)
q.Enqueue(2)
val, _ := q.Dequeue() // returns 1
```

### Sorting Algorithms

The `sorting` package provides several classic sorting algorithms:

- **Selection Sort**: O(n²) comparison-based sort
- **Insertion Sort**: O(n²) adaptive sort, efficient for small arrays
- **Merge Sort**: O(n log n) divide-and-conquer stable sort
- **Quick Sort**: O(n log n) average case, in-place partitioning sort
- **Heap Sort**: O(n log n) in-place sort using binary heap
- **Quick Select**: O(n) average case for finding kth element

```go
import "github.com/jjmrocha/go-algo/sorting"

arr := []int{64, 25, 12, 22, 11}
sorting.QuickSort(arr) // arr is now [11, 12, 22, 25, 64]

// Find the kth smallest element
kth := sorting.QuickSelect(arr, 2) // returns 3rd smallest element
```

### Priority Queue

Both minimum and maximum priority queue implementations using binary heaps:

```go
import "github.com/jjmrocha/go-algo/priorityqueue"

// Max Priority Queue
maxPQ := priorityqueue.NewMaxPQ[int]()
maxPQ.Insert(5)
maxPQ.Insert(3)
max, _ := maxPQ.DelMax() // returns 5

// Min Priority Queue
minPQ := priorityqueue.NewMinPQ[int]()
minPQ.Insert(5)
minPQ.Insert(3)
min, _ := minPQ.DelMin() // returns 3
```

### Dictionaries (Associative Arrays)

The `dict` package provides a Binary Search Tree implementation:

```go
import "github.com/jjmrocha/go-algo/dict"

bst := dict.NewBST[int, string]()
bst.Put(5, "five")
bst.Put(3, "three")
val, ok := bst.Get(5) // returns "five", true
bst.Delete(3)
```

### Tree Structures

The `trees` package includes several advanced tree structures:

#### Red-Black Binary Search Tree

Self-balancing binary search tree with O(log n) operations:

```go
import "github.com/jjmrocha/go-algo/trees"

rbt := trees.NewRedBlackBST[int, string]()
rbt.Put(5, "five")
rbt.Put(3, "three")
val, ok := rbt.Get(5) // returns "five", true
```

#### B-Tree

Multi-way search tree optimized for systems with large blocks of data:

```go
bt := trees.NewBTree[int, string](3) // minimum degree = 3
bt.Insert(5, "five")
bt.Insert(3, "three")
val, ok := bt.Search(5) // returns "five", true
```

#### kd-Tree

2D tree for efficient spatial searches:

```go
kd := trees.NewKDTree()
kd.Insert(trees.Point2D{X: 1.0, Y: 2.0})
kd.Insert(trees.Point2D{X: 3.0, Y: 4.0})

// Find nearest point
nearest := kd.Nearest(trees.Point2D{X: 2.9, Y: 4.1})

// Range search
points := kd.RangeSearch(0.0, 0.0, 5.0, 5.0)
```

### Hash Tables

The `hashtable` package provides two hash table implementations:

#### Separate Chaining

Hash table using linked lists to handle collisions:

```go
import "github.com/jjmrocha/go-algo/hashtable"

ht := hashtable.NewSeparateChainingHashTable[int, string](10)
ht.Put(1, "one")
ht.Put(2, "two")
val, ok := ht.Get(1) // returns "one", true
ht.Delete(2)
```

#### Linear Probing

Hash table using open addressing with linear probing:

```go
ht := hashtable.NewLinearProbingHashTable[int, string](16)
ht.Put(1, "one")
ht.Put(2, "two")
val, ok := ht.Get(1) // returns "one", true
```

## Installation

```bash
go get github.com/jjmrocha/go-algo
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test ./... -v
```

## Requirements

- Go 1.25 or later (uses generics)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## References

This implementation is based on the algorithms taught in:
- **Algorithms, Part I** - Princeton University on Coursera
- Instructors: Robert Sedgewick and Kevin Wayne