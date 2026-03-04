# go-algo

A collection of generic algorithms and data structures implemented in Go, designed to be reused across projects.
Each package is independently importable and carries no cross-package dependencies (except `cache`, which uses `future` internally).

Requires **Go 1.23+** (uses `iter.Seq` / range-over-func).

## Installation

```bash
go get github.com/jjmrocha/go-algo
```

## Packages

### `stack` — Generic LIFO stack

```go
import "github.com/jjmrocha/go-algo/stack"

s := stack.New[int]()
s.Push(1)
s.Push(2)
v, ok := s.Pop()   // 2, true
v, ok  = s.Peek()  // 1, true (no removal)
```

`SyncStack[T]` provides the same API safe for concurrent use.

---

### `queue` — Generic FIFO queue

```go
import "github.com/jjmrocha/go-algo/queue"

q := queue.New[string]()
q.Enqueue("a")
q.Enqueue("b")
v, ok := q.Dequeue() // "a", true
```

`SyncQueue[T]` provides the same API safe for concurrent use.

---

### `sets` — Generic set with set-algebra operations

```go
import "github.com/jjmrocha/go-algo/sets"

a := sets.Of([]int{1, 2, 3})
b := sets.Of([]int{2, 3, 4})

a.Add(5)
a.Remove(1)
a.Contains(2) // true

a.Union(b)        // {2, 3, 4, 5}
a.Intersection(b) // {2, 3}
a.Difference(b)   // {5}

for v := range a.Values() { /* ... */ }
```

A nil `Set` is safe for all read operations (`Contains`, `Len`, `ToSlice`).

---

### `union` — Weighted quick-union (disjoint sets)

```go
import "github.com/jjmrocha/go-algo/union"

u := union.New(10) // 10 elements, each its own set

u.Union(0, 1)
u.Union(1, 2)

connected, _ := u.Connected(0, 2) // true
root, _       := u.Find(2)        // root of the set containing 2
```

Uses union by weight and half-path compression for near-constant time per operation.

---

### `functional` — Higher-order functions for slices and iterators

Every operation comes in two variants: a slice variant and a lazy iterator variant (`Seq` suffix).

```go
import "github.com/jjmrocha/go-algo/functional"

nums := []int{1, 2, 3, 4, 5}

functional.Filter(nums, func(v int) bool { return v%2 == 0 }) // [2, 4]
functional.Map(nums, func(v int) int { return v * 10 })        // [10, 20, 30, 40, 50]
functional.Fold(nums, 0, func(acc, v int) int { return acc + v }) // 15
functional.Distinct([]int{1, 2, 1, 3})                        // [1, 2, 3]
functional.Find(nums, func(v int) bool { return v > 3 })      // 4, true

matching, rest := functional.Partition(nums, func(v int) bool { return v%2 == 0 })
// matching=[2,4], rest=[1,3,5]

functional.GroupBy(nums, func(v int) string {
    if v%2 == 0 { return "even" }
    return "odd"
}) // map[even:[2 4] odd:[1 3 5]]
```

Iterator variants integrate with Go's `iter.Seq[T]` for lazy, composable pipelines:

```go
import "slices"

seq := functional.FilterSeq(slices.Values(nums), func(v int) bool { return v%2 == 0 })
seq  = functional.MapSeq(seq, func(v int) int { return v * 10 })
result := slices.Collect(seq) // [20, 40]
```

---

### `future` — Generic async computation

```go
import "github.com/jjmrocha/go-algo/future"

f := future.Async(func() (int, error) {
    return expensiveComputation(), nil
})

// Multiple goroutines can Await the same Future concurrently.
result, err := f.Await(ctx)

// Or with an explicit timeout:
result, err = f.AwaitWithTimeout(5 * time.Second)
```

`AsyncWithContext` forwards a context to the computation function.

---

### `cache` — LRU cache

#### Direct cache — `LRUCache[K, V]`

```go
import "github.com/jjmrocha/go-algo/cache"

c, err := cache.NewLRUCache[string, int](100)

c.Put("hits", 42)
v, ok := c.Get("hits") // 42, true — promotes to MRU

c.Exists("hits") // true — does NOT affect LRU order
c.Delete("hits")
```

#### Auto-loading cache — `LRUProvider[K, V]`

Fetches missing values from a `Provider` function, deduplicating concurrent requests for the same key.

```go
provider := func(key string) (string, error) {
    return fetchFromDB(key)
}

lp, err := cache.NewLRUWithProvider[string, string](100, provider)

// Cache miss → provider called once, result shared with concurrent callers.
value, err := lp.Get(ctx, "key")

// Force re-fetch on next Get.
lp.Invalidate("key")
```

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
