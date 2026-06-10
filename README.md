# go-algo

A collection of generic algorithms and data structures for Go, designed to be reused across projects.
Each package is independently importable, and the runtime has no third-party dependencies.

Requires **Go 1.25+** (uses `iter.Seq` / range-over-func, the `min` builtin, and `maps.Copy`).

## Installation

```bash
go get github.com/jjmrocha/go-algo
```

## Packages at a glance

| Package | Category | Summary |
|---------|----------|---------|
| [`stack`](#stack) | Containers | LIFO stack, plus a concurrent `SyncStack` |
| [`queue`](#queue) | Containers | FIFO, bounded-blocking, and priority queues (plus sync variants) |
| [`deque`](#deque) | Containers | Double-ended queue backed by a ring buffer |
| [`sets`](#sets) | Membership & counting | Set with union / intersection / difference |
| [`bag`](#bag) | Membership & counting | Multiset tracking a count per element |
| [`treemap`](#treemap) | Ordered & disjoint-set | Sorted key-value map (Left-Leaning Red-Black tree) |
| [`unionfind`](#unionfind) | Ordered & disjoint-set | Disjoint-set / union-find |
| [`sorting`](#sorting) | Algorithms | Insertion, Shell, Merge, Quick, and Shuffle |
| [`fn`](#fn) | Algorithms | Map / Filter / Fold / … for slices and iterators |
| [`future`](#future) | Concurrency | Async result with await, timeout, and cancellation |
| [`singleflight`](#singleflight) | Concurrency | Duplicate-call suppression |
| [`cache`](#cache) | Concurrency | LRU cache and auto-loading provider |

**Internal dependencies:** most packages are standalone. `treemap` uses `stack`; `singleflight` uses `future`; `cache` uses `future` and `singleflight`.

### The comparator contract

The ordered types (`treemap`, `queue.PriorityQueue`) and the `sorting` algorithms take a
comparator `func(a, b T) int` that returns a **negative number, zero, or a positive number** —
the same convention as the standard library's [`cmp.Compare`](https://pkg.go.dev/cmp#Compare) and
[`slices.SortFunc`](https://pkg.go.dev/slices#SortFunc). Pass `cmp.Compare[T]` for natural ordering,
or supply your own function for a custom order.

---

## Containers

### `stack`

Generic LIFO stack backed by a singly-linked list.

```go
import "github.com/jjmrocha/go-algo/stack"

s := stack.New[int]()
s.Push(1)
s.Push(2)
v, ok := s.Pop()  // 2, true
v, ok  = s.Peek() // 1, true (no removal)

for v := range s.Drain() { /* pops top→bottom; the stack is emptied */ }
```

`SyncStack[T]` provides the same API, safe for concurrent use.

### `queue`

Generic FIFO queue backed by a singly-linked list.

```go
import "github.com/jjmrocha/go-algo/queue"

q := queue.New[string]()
q.Enqueue("a")
q.Enqueue("b")
v, ok := q.Dequeue() // "a", true

for v := range q.Drain() { /* dequeues front→back; the queue is emptied */ }
```

`SyncQueue[T]` provides the same API, safe for concurrent use.

**`BlockingQueue[T]`** is a bounded, channel-backed variant. `Enqueue` blocks when full;
`Dequeue` blocks when empty — suitable for producer/consumer pipelines.

```go
bq, err := queue.NewBlockingQueue[int](16) // capacity 16
bq.Enqueue(42)                             // blocks if full
v := bq.Dequeue()                          // blocks if empty
```

**`PriorityQueue[T]`** is a min-heap ordered by a comparator. Elements are dequeued in priority
order, not insertion order.

```go
import "cmp"

pq := queue.NewPriorityQueue[int](cmp.Compare[int])
pq.Enqueue(5)
pq.Enqueue(1)
pq.Enqueue(3)

v, ok := pq.Peek()    // 1, true — minimum, no removal
v, ok  = pq.Dequeue() // 1, true — removes minimum

// Custom initial capacity (returns an error if ≤ 0).
pq2, err := queue.NewPriorityQueueWithCap[int](64, cmp.Compare[int])

for v := range pq2.Drain() { /* drains in priority order */ }
```

`SyncPriorityQueue[T]` provides the same API, safe for concurrent use.

### `deque`

Generic double-ended queue backed by a dynamically resized ring buffer, with O(1) amortised push
and pop at both ends.

```go
import "github.com/jjmrocha/go-algo/deque"

d := deque.New[int]()
d.PushBack(1)
d.PushBack(2)
d.PushFront(0)        // logical order front→back: 0, 1, 2

v, ok := d.PopFront() // 0, true
v, ok  = d.PopBack()  // 2, true
v, ok  = d.PeekFront() // 1, true (no removal)

for v := range d.Values() { /* front→back; the deque is NOT modified */ }
```

`SyncDeque[T]` provides the same API, safe for concurrent use.

---

## Membership & counting

### `sets`

Generic set of unique elements with set-algebra operations. O(1) amortised `Add`, `Remove`, and
`Contains`.

```go
import "github.com/jjmrocha/go-algo/sets"

a := sets.Of([]int{1, 2, 3})
b := sets.Of([]int{2, 3, 4})

a.Contains(2)     // true
a.Union(b)        // {1, 2, 3, 4}
a.Intersection(b) // {2, 3}
a.Difference(b)   // {1}

for v := range a.Values() { /* unspecified order; not modified */ }
```

A nil `Set` is safe for read operations (`Contains`, `Len`, `ToSlice`) and for `Remove` (a no-op);
`Add` on a nil `Set` panics.

### `bag`

Generic multiset that tracks how many times each element has been added, allowing duplicates.

```go
import "github.com/jjmrocha/go-algo/bag"

b := bag.New[string]()
b.Add("a", "b", "a")  // "a" added twice

b.Contains("a") // true
b.Count("a")    // 2
b.Len()         // 3 — total items including duplicates
b.Unique()      // ["a", "b"] — one entry per distinct element (unspecified order)
```

`Add`, `Remove`, `Contains`, and `Count` are O(1) amortised; `Len`, `ToSlice`, `Unique`, and
`Values` are O(n). Set-algebra on counts: `Union` sums counts, `Intersection` takes the minimums.
A nil `Bag` is safe for read operations (`Contains`, `Len`, `Count`, `Empty`).

---

## Ordered & disjoint-set

### `treemap`

Ordered key-value store backed by a Left-Leaning Red-Black BST. Keys are kept in sorted order.

```go
import (
    "cmp"

    "github.com/jjmrocha/go-algo/treemap"
)

m := treemap.New[int, string](cmp.Compare[int])
m.Put(2, "b")
m.Put(1, "a")
m.Put(3, "c")

v, ok := m.Get(1) // "a", true
m.ToList()        // ["a", "b", "c"] — values in ascending key order

k, ok := m.Min()  // 1, true — smallest key
m.Rank(2)         // 1 — number of keys strictly less than 2
k, ok = m.Select(0) // 1, true — key of rank 0 (minimum)
```

`Get`, `Put`, `Contains`, `Delete`, `Min`, and `Max` run in O(log n). `Rank` and `Select` run in
O(n) because nodes do not cache subtree sizes.

### `unionfind`

Union-find (disjoint-set) structure using weighted quick-union with half-path compression, giving
near-constant time per operation.

```go
import "github.com/jjmrocha/go-algo/unionfind"

u := unionfind.New(10) // 10 elements, each its own set

u.Union(0, 1)
u.Union(1, 2)

connected, _ := u.Connected(0, 2) // true
root, _       := u.Find(2)        // root of the set containing 2
u.Sets()                          // 8 — number of disjoint sets
```

---

## Algorithms

### `sorting`

Generic in-place sorting and shuffling. All sort functions take a comparator (see
[The comparator contract](#the-comparator-contract)).

```go
import (
    "cmp"

    "github.com/jjmrocha/go-algo/sorting"
)

arr := []int{4, 3, 6, 1, 5, 2}

sorting.Insertion(arr, cmp.Compare[int]) // stable, O(n²) — small or nearly-sorted inputs
sorting.Shell(arr, cmp.Compare[int])     // O(n log² n) — Ciura gap sequence, fast in practice
sorting.Merge(arr, cmp.Compare[int])     // stable, O(n log n) — single auxiliary buffer
sorting.Quick(arr, cmp.Compare[int])     // O(n log n) avg — 3-way partition, middle-index pivot

// Custom order: any func(a, b T) int returning <0 / 0 / >0.
desc := func(a, b int) int { return cmp.Compare(b, a) }
sorting.Quick(arr, desc)
```

`Shuffle` randomly permutes a slice using a cryptographically secure source:

```go
err := sorting.Shuffle(arr) // Fisher-Yates via crypto/rand
```

### `fn`

Higher-order functions for slices and iterators.

```go
import "github.com/jjmrocha/go-algo/fn"

nums := []int{1, 2, 3, 4, 5}

fn.Filter(nums, func(v int) bool { return v%2 == 0 }) // [2, 4]
fn.Map(nums, func(v int) int { return v * 10 })        // [10, 20, 30, 40, 50]
fn.Fold(nums, 0, func(acc, v int) int { return acc + v }) // 15
fn.Distinct([]int{1, 2, 1, 3})                         // [1, 2, 3]

matching, rest := fn.Partition(nums, func(v int) bool { return v%2 == 0 })
// matching=[2,4], rest=[1,3,5]
```

Most operations also have a lazy iterator variant (`Seq` suffix) that composes with Go's
`iter.Seq[T]`:

```go
import "slices"

seq := fn.FilterSeq(slices.Values(nums), func(v int) bool { return v%2 == 0 })
seq  = fn.MapSeq(seq, func(v int) int { return v * 10 })
result := slices.Collect(seq) // [20, 40]
```

`Partition` and `Zip` are slice-only; `GroupBySeq` materialises its groups before yielding (grouping
is inherently eager).

---

## Concurrency

### `future`

Generic async computation. A `Future[T]` holds a value that becomes available later.

```go
import "github.com/jjmrocha/go-algo/future"

f := future.Async(func() (int, error) {
    return expensiveComputation(), nil
})

result, err := f.Await()                       // block until done
result, err  = f.AwaitWithContext(ctx)         // block with cancellation / deadline
result, err  = f.AwaitWithTimeout(5 * time.Second) // block with a timeout
```

Multiple goroutines can await the same Future concurrently — all receive the same result.
`AsyncWithContext` forwards a context to the computation.

### `singleflight`

Deduplicates concurrent calls for the same key, preventing cache stampedes. When multiple goroutines
request the same key simultaneously, only one provider invocation runs; all callers share the result.

```go
import "github.com/jjmrocha/go-algo/singleflight"

sf := singleflight.New[string, User]()

f := sf.Do("user:42", func() (User, error) {
    return fetchUserFromDB(42)
})
user, err := f.Await()
```

Completed computations are not retained: a call after the Future resolves starts a fresh invocation.

### `cache`

LRU cache. Both cache types are configured with functional options; `WithCapacity` is required and
`WithTTL` is optional.

```go
import (
    "time"

    "github.com/jjmrocha/go-algo/cache"
)

c, err := cache.NewLRUCache[string, int](
    cache.WithCapacity(100),
    cache.WithTTL(5 * time.Minute), // optional — omit for no expiry
)

c.Put("hits", 42)
v, ok := c.Get("hits") // 42, true — promotes to most-recently-used
c.Exists("hits")       // true — does NOT affect LRU order
```

**`LRUProvider[K, V]`** fetches missing values from a `Provider` function, deduplicating concurrent
requests for the same key (via `singleflight`):

```go
lp, err := cache.NewLRUWithProvider(
    func(key string) (string, error) { return fetchFromDB(key) },
    cache.WithCapacity(100),
    cache.WithTTL(5 * time.Minute), // optional
)

value, err := lp.Get("key")                 // blocks until the provider returns
value, err  = lp.GetWithContext(ctx, "key") // context-aware
lp.Invalidate("key")                        // force re-fetch on next Get
```

---

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
