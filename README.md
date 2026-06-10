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

Generic LIFO stack backed by a singly-linked list. · [API ↓](#stack-api)

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

Generic FIFO queue backed by a singly-linked list. · [API ↓](#queue-api)

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
and pop at both ends. · [API ↓](#deque-api)

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
`Contains`. · [API ↓](#sets-api)

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

Generic multiset that tracks how many times each element has been added, allowing duplicates. · [API ↓](#bag-api)

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

Ordered key-value store backed by a Left-Leaning Red-Black BST. Keys are kept in sorted order. · [API ↓](#treemap-api)

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
near-constant time per operation. · [API ↓](#unionfind-api)

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
[The comparator contract](#the-comparator-contract)). · [API ↓](#sorting-api)

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

Higher-order functions for slices and iterators. · [API ↓](#fn-api)

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

Generic async computation. A `Future[T]` holds a value that becomes available later. · [API ↓](#future-api)

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
request the same key simultaneously, only one provider invocation runs; all callers share the result. · [API ↓](#singleflight-api)

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
`WithTTL` is optional. · [API ↓](#cache-api)

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

## API reference

Complete exported surface of each package. For full documentation see
[pkg.go.dev/github.com/jjmrocha/go-algo](https://pkg.go.dev/github.com/jjmrocha/go-algo).

### stack API

Constructors: `New[T]() *Stack[T]`, `NewSyncStack[T]() *SyncStack[T]`.

| Method | Signature | Description |
|--------|-----------|-------------|
| Push | `Push(data T)` | Add an element to the top. |
| Pop | `Pop() (T, bool)` | Remove and return the top; `false` if empty. |
| Peek | `Peek() (T, bool)` | Return the top without removing; `false` if empty. |
| Len | `Len() int` | Number of elements. |
| Empty | `Empty() bool` | Whether the stack is empty. |
| Drain | `Drain() iter.Seq[T]` | Iterator popping top→bottom (empties the stack). |

`SyncStack[T]` exposes the same methods, safe for concurrent use.

### queue API

Constructors: `New[T]() *Queue[T]`, `NewSyncQueue[T]() *SyncQueue[T]`,
`NewBlockingQueue[T](capacity int) (*BlockingQueue[T], error)`,
`NewPriorityQueue[T](cmp func(a, b T) int) *PriorityQueue[T]`,
`NewPriorityQueueWithCap[T](initialCap int, cmp func(a, b T) int) (*PriorityQueue[T], error)`,
`NewSyncPriorityQueue[T](cmp func(a, b T) int) *SyncPriorityQueue[T]`,
`NewSyncPriorityQueueWithCap[T](initialCap int, cmp func(a, b T) int) (*SyncPriorityQueue[T], error)`.
Errors: `ErrCapacityTooSmall`.

`Queue[T]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| Enqueue | `Enqueue(data T)` | Add to the back. |
| Dequeue | `Dequeue() (T, bool)` | Remove and return the front; `false` if empty. |
| Len | `Len() int` | Number of elements. |
| Empty | `Empty() bool` | Whether empty. |
| Drain | `Drain() iter.Seq[T]` | Iterator dequeuing front→back (empties the queue). |

`BlockingQueue[T]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| Enqueue | `Enqueue(v T)` | Add to the back; blocks while full. |
| Dequeue | `Dequeue() T` | Remove and return the front; blocks while empty. |
| Len | `Len() int` | Current number of elements. |
| Empty | `Empty() bool` | Whether empty. |
| Cap | `Cap() int` | Maximum capacity. |
| Full | `Full() bool` | Whether at capacity. |

`PriorityQueue[T]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| Enqueue | `Enqueue(data T)` | Insert, restoring heap order. |
| Peek | `Peek() (T, bool)` | Minimum without removing; `false` if empty. |
| Dequeue | `Dequeue() (T, bool)` | Remove and return the minimum; `false` if empty. |
| Len | `Len() int` | Number of elements. |
| Empty | `Empty() bool` | Whether empty. |
| Drain | `Drain() iter.Seq[T]` | Iterator draining in priority order. |

`SyncQueue[T]` mirrors `Queue[T]`; `SyncPriorityQueue[T]` mirrors `PriorityQueue[T]` — same methods, safe for concurrent use.

### deque API

Constructors: `New[T]() *Deque[T]`, `NewWithCap[T](initialCap int) (*Deque[T], error)`,
`NewSyncDeque[T]() *SyncDeque[T]`. Errors: `ErrCapacityTooSmall`.

`Deque[T]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| PushFront | `PushFront(data T)` | Add to the front. |
| PushBack | `PushBack(data T)` | Add to the back. |
| PopFront | `PopFront() (T, bool)` | Remove and return the front; `false` if empty. |
| PopBack | `PopBack() (T, bool)` | Remove and return the back; `false` if empty. |
| PeekFront | `PeekFront() (T, bool)` | Front without removing; `false` if empty. |
| PeekBack | `PeekBack() (T, bool)` | Back without removing; `false` if empty. |
| Len | `Len() int` | Number of elements. |
| Empty | `Empty() bool` | Whether empty. |
| Values | `Values() iter.Seq[T]` | Front→back iterator (non-destructive). |
| ToSlice | `ToSlice() []T` | Copy of all elements, front→back. |

`SyncDeque[T]` exposes the same methods, safe for concurrent use.

### sets API

Constructors: `New[T comparable](items ...T) Set[T]`, `Of[T comparable](items []T) Set[T]`
(`Set[T]` is `map[T]struct{}`).

| Method | Signature | Description |
|--------|-----------|-------------|
| Add | `Add(items ...T)` | Insert elements (panics on a nil set). |
| Remove | `Remove(items ...T)` | Delete elements (nil-safe no-op). |
| Contains | `Contains(value T) bool` | Membership test. |
| Len | `Len() int` | Number of elements. |
| Empty | `Empty() bool` | Whether empty. |
| ToSlice | `ToSlice() []T` | Elements as a slice (unspecified order). |
| Values | `Values() iter.Seq[T]` | Iterator over elements (non-destructive). |
| Union | `Union(o Set[T]) Set[T]` | New set: elements in either. |
| Intersection | `Intersection(o Set[T]) Set[T]` | New set: elements in both. |
| Difference | `Difference(o Set[T]) Set[T]` | New set: in this but not `o`. |
| String | `String() string` | `set{...}` representation. |

### bag API

Constructors: `New[T comparable](items ...T) Bag[T]`, `Of[T comparable](items []T) Bag[T]`
(`Bag[T]` is `map[T]int`).

| Method | Signature | Description |
|--------|-----------|-------------|
| Add | `Add(items ...T)` | Increment counts (panics on a nil bag). |
| Remove | `Remove(items ...T)` | Decrement counts; deletes at zero. |
| Clear | `Clear()` | Remove all items. |
| Contains | `Contains(value T) bool` | Whether present at least once. |
| Count | `Count(value T) int` | Occurrence count of a value. |
| Len | `Len() int` | Total items including duplicates (O(n)). |
| Empty | `Empty() bool` | Whether empty. |
| ToSlice | `ToSlice() []T` | Each item repeated by its count. |
| Unique | `Unique() []T` | One entry per distinct item. |
| Values | `Values() iter.Seq[T]` | Iterator yielding each item by its count. |
| Union | `Union(o Bag[T]) Bag[T]` | New bag: counts summed. |
| Intersection | `Intersection(o Bag[T]) Bag[T]` | New bag: minimum counts. |
| String | `String() string` | `bag{...}` representation. |

### treemap API

Constructor: `New[K, V any](cmp func(a, b K) int) *Map[K, V]`.

| Method | Signature | Description |
|--------|-----------|-------------|
| Put | `Put(key K, value V)` | Insert or replace. |
| Get | `Get(key K) (V, bool)` | Value for key; `false` if absent. |
| Contains | `Contains(key K) bool` | Whether key is present. |
| Delete | `Delete(key K) bool` | Remove key; `false` if absent. |
| Len | `Len() int` | Number of pairs. |
| Empty | `Empty() bool` | Whether empty. |
| Min | `Min() (K, bool)` | Smallest key; `false` if empty. |
| Max | `Max() (K, bool)` | Largest key; `false` if empty. |
| Rank | `Rank(key K) int` | Count of keys strictly less than key (O(n)). |
| Select | `Select(r int) (K, bool)` | Key of rank `r` (O(n)); `false` if out of range. |
| ToList | `ToList() []V` | Values in ascending key order. |

### unionfind API

Constructor: `New(size int) *UnionFind`. Errors: `ErrIndexOutOfRange`.

| Method | Signature | Description |
|--------|-----------|-------------|
| Union | `Union(p, q int) error` | Merge the sets of `p` and `q`. |
| Find | `Find(p int) (int, error)` | Root of `p`'s set (path-compressing). |
| Connected | `Connected(p, q int) (bool, error)` | Whether `p` and `q` share a set. |
| Sets | `Sets() int` | Number of disjoint sets. |
| String | `String() string` | Parent-link representation. |

### sorting API

| Function | Signature | Description |
|----------|-----------|-------------|
| Insertion | `Insertion[T]([]T, Comparator[T])` | In-place insertion sort (stable, O(n²)). |
| Shell | `Shell[T]([]T, Comparator[T])` | In-place Shell sort (Ciura gaps). |
| Merge | `Merge[T]([]T, Comparator[T])` | Stable merge sort (O(n log n)). |
| Quick | `Quick[T]([]T, Comparator[T])` | 3-way quicksort (middle-index pivot). |
| Shuffle | `Shuffle[T]([]T) error` | Fisher–Yates via crypto/rand. |
| ShuffleWithRandom | `ShuffleWithRandom[T]([]T, RandomNextInt) error` | Shuffle with an injectable random source. |
| Swap | `Swap[T]([]T, i, j int)` | Exchange two elements. |

Types: `Comparator[T] = func(a, b T) int`, `RandomNextInt = func(int) (int, error)`.

### fn API

Each slice function has a lazy `Seq` twin unless noted.

| Function (+ Seq twin) | Signature | Description |
|-----------------------|-----------|-------------|
| Map / MapSeq | `Map[T, U]([]T, func(T) U) []U` | Transform each element. |
| Filter / FilterSeq | `Filter[T]([]T, func(T) bool) []T` | Keep elements matching a predicate. |
| Fold / FoldSeq | `Fold[T, U]([]T, U, func(U, T) U) U` | Reduce to a single value. |
| Distinct / DistinctSeq | `Distinct[T comparable]([]T) []T` | Remove duplicates (first occurrence wins). |
| Find / FindSeq | `Find[T]([]T, func(T) bool) (T, bool)` | First match; `false` if none. |
| Any / AnySeq | `Any[T]([]T, func(T) bool) bool` | At least one element matches. |
| All / AllSeq | `All[T]([]T, func(T) bool) bool` | All match (vacuously true if empty). |
| ForEach / ForEachSeq | `ForEach[T]([]T, func(T))` | Run a side effect per element. |
| GroupBy / GroupBySeq | `GroupBy[T, K comparable]([]T, func(T) K) map[K][]T` | Bucket by key (the Seq variant is eager). |
| Partition | `Partition[T]([]T, func(T) bool) ([]T, []T)` | Split matching / non-matching (slice-only). |
| Zip | `Zip[T, U, V]([]T, []U, func(T, U) V) []V` | Combine two slices pairwise (slice-only). |

### future API

Constructors: `Async[T](provider func() (T, error)) *Future[T]`,
`AsyncWithContext[T](ctx context.Context, provider func(context.Context) (T, error)) *Future[T]`.

| Method | Signature | Description |
|--------|-----------|-------------|
| Await | `Await() (T, error)` | Block until resolved. |
| AwaitWithContext | `AwaitWithContext(ctx context.Context) (T, error)` | Block until resolved or ctx is done. |
| AwaitWithTimeout | `AwaitWithTimeout(d time.Duration) (T, error)` | Block until resolved or the timeout elapses. |
| Done | `Done() bool` | Whether the computation has resolved. |

### singleflight API

Constructor: `New[K comparable, V any]() *SingleFlight[K, V]`.

| Method | Signature | Description |
|--------|-----------|-------------|
| Do | `Do(key K, provider func() (V, error)) *future.Future[V]` | Share one in-flight call per key. |

### cache API

Constructors: `NewLRUCache[K, V](opts ...Option) (*LRUCache[K, V], error)`,
`NewLRUWithProvider[K, V](provider Provider[K, V], opts ...Option) (*LRUProvider[K, V], error)`.
Options: `WithCapacity(n int) Option` (required), `WithTTL(d time.Duration) Option`.
Types: `Option`, `Provider[K, V] = func(key K) (V, error)`.
Errors: `ErrInvalidCapacity`, `ErrInvalidTTL`, `ErrNilProvider`.

`LRUCache[K, V]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| Get | `Get(key K) (V, bool)` | Value for key; promotes to most-recently-used. |
| Put | `Put(key K, value V)` | Insert/update; evicts the LRU entry at capacity. |
| Delete | `Delete(key K)` | Remove key (no-op if absent). |
| Exists | `Exists(key K) bool` | Presence test; does not affect LRU order. |
| Len | `Len() int` | Number of entries. |
| Cap | `Cap() int` | Maximum entries. |

`LRUProvider[K, V]`:

| Method | Signature | Description |
|--------|-----------|-------------|
| Get | `Get(key K) (V, error)` | Cached value, loading via the provider on a miss (blocks). |
| GetWithContext | `GetWithContext(ctx context.Context, key K) (V, error)` | Context-aware `Get`. |
| Exists | `Exists(key K) bool` | Presence test; does not affect LRU order. |
| Invalidate | `Invalidate(key K)` | Drop key; the next `Get` re-fetches. |
| Len | `Len() int` | Number of entries. |
| Cap | `Cap() int` | Maximum entries. |

---

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
