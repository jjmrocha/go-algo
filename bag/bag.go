// Package bag provides a generic multiset (bag) — an unordered collection
// that, unlike a set, allows duplicate elements and tracks how many times
// each element has been added.
//
// A [Bag] is backed by a map[T]int and provides O(1) amortised [Bag.Add],
// [Bag.Remove], and [Bag.Contains] operations. Bags are not safe for
// concurrent use without external synchronisation.
package bag

import (
	"fmt"
	"iter"
	"maps"
	"strings"
)

// Bag is a generic multiset that maps each element to its occurrence count.
//
// The zero value of Bag is nil. A nil Bag is safe for read operations
// ([Bag.Contains], [Bag.Len], [Bag.Count], [Bag.Empty]), but calling
// [Bag.Add] on a nil Bag panics. Use [New] or [Of] to obtain a
// ready-to-use Bag.
//
// Because Bag is a map type, assigning one Bag variable to another creates
// an alias: both variables refer to the same underlying collection.
type Bag[T comparable] map[T]int

// New returns an empty Bag, optionally pre-populated with items.
func New[T comparable](items ...T) Bag[T] {
	return Of(items)
}

// Of creates a Bag from a slice of items, recording the count of each.
func Of[T comparable](items []T) Bag[T] {
	s := make(Bag[T], len(items))
	s.Add(items...)
	return s
}

// Add increments the count of each item in b by one.
// Calling Add on a nil Bag panics; use [New] or [Of] to create a
// non-nil Bag first.
func (b Bag[T]) Add(items ...T) {
	for _, item := range items {
		b[item]++
	}
}

// Remove decrements the count of each item in b by one, deleting items
// whose count reaches zero. Removing an absent item is a no-op.
func (b Bag[T]) Remove(items ...T) {
	for _, item := range items {
		if count, exists := b[item]; exists {
			if count > 1 {
				b[item] = count - 1
			} else {
				delete(b, item)
			}
		}
	}
}

// Clear removes all items from b.
func (b Bag[T]) Clear() {
	for k := range b {
		delete(b, k)
	}
}

// Contains reports whether value is present in b at least once.
// Contains on a nil Bag always returns false.
func (b Bag[T]) Contains(value T) bool {
	_, exists := b[value]
	return exists
}

// Len returns the total number of items in b, counting duplicates.
// Len on a nil Bag returns 0.
func (b Bag[T]) Len() int {
	total := 0

	for _, count := range b {
		total += count
	}

	return total
}

// Count returns the number of times value appears in b.
// Count on a nil Bag always returns 0.
func (b Bag[T]) Count(value T) int {
	return b[value]
}

// Empty reports whether b contains no items.
// Empty on a nil Bag returns true.
func (b Bag[T]) Empty() bool {
	return len(b) == 0
}

// ToSlice returns a new slice containing all items in b, with each item
// repeated according to its count. The order of items is unspecified.
func (b Bag[T]) ToSlice() []T {
	var result []T

	for v, count := range b {
		for range count {
			result = append(result, v)
		}
	}

	return result
}

// Unique returns a new slice containing each distinct item in b exactly
// once. The order of items is unspecified.
func (b Bag[T]) Unique() []T {
	result := make([]T, 0, len(b))

	for v := range b {
		result = append(result, v)
	}

	return result
}

// Values returns an iterator that yields each item in b, with each item
// repeated according to its count. The order of items is unspecified.
func (b Bag[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, count := range b {
			for range count {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// String returns a human-readable representation of b in the form
// bag{v1: c1, v2: c2, ...}. The order of entries is unspecified.
// A nil Bag returns "bag(nil)" and an empty Bag returns "bag{}".
func (b Bag[T]) String() string {
	if b == nil {
		return "bag(nil)"
	}

	var sb strings.Builder
	sb.WriteString("bag{")
	first := true

	for v, count := range b {
		if !first {
			sb.WriteString(", ")
		}

		_, _ = fmt.Fprintf(&sb, "%v: %d", v, count)
		first = false
	}

	sb.WriteString("}")

	return sb.String()
}

// Union returns a new Bag containing all items from b and o, with counts
// summed. A nil operand is treated as an empty bag.
func (b Bag[T]) Union(o Bag[T]) Bag[T] {
	result := make(Bag[T], len(b))

	maps.Copy(result, b)

	for v, count := range o {
		result[v] += count
	}

	return result
}

// Intersection returns a new Bag containing only the items present in both
// b and o, with counts set to the minimum of the two. A nil operand is
// treated as an empty bag.
func (b Bag[T]) Intersection(o Bag[T]) Bag[T] {
	result := make(Bag[T])

	for v, count := range b {
		if oCount := o.Count(v); oCount > 0 {
			result[v] = min(count, oCount)
		}
	}

	return result
}
