// Package sets provides a generic, unordered collection of unique elements.
//
// A [Set] is backed by a map and provides O(1) amortised [Set.Add],
// [Set.Remove], and [Set.Contains] operations. Sets are not safe for
// concurrent use without external synchronisation.
package sets

import (
	"fmt"
	"iter"
	"strings"
)

// Set is an unordered collection of unique elements of type T.
//
// The zero value of Set is nil. A nil Set is safe for read operations
// ([Set.Contains], [Set.Len], [Set.Remove], [Set.ToSlice]), but calling
// [Set.Add] on a nil Set panics. Use [New] or [Of] to obtain a
// ready-to-use Set.
//
// Because Set is a map type, assigning one Set variable to another creates
// an alias: both variables refer to the same underlying collection.
type Set[T comparable] map[T]struct{}

// New returns an empty, initialized Set ready for use.
func New[T comparable](items ...T) Set[T] {
	return Of(items)
}

// Of creates a Set containing the elements of items, silently
// discarding any duplicates.
func Of[T comparable](items []T) Set[T] {
	s := make(Set[T], len(items))
	s.Add(items...)
	return s
}

// Add inserts items into s. Adding an element that is already present is a
// no-op. Calling Add on a nil Set panics; use [New] or [Of] to
// create a non-nil Set first.
func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Remove deletes items from s. Removing an absent element is a no-op.
// Calling Remove on a nil Set is safe.
func (s Set[T]) Remove(items ...T) {
	for _, item := range items {
		delete(s, item)
	}
}

// Contains reports whether value is an element of s.
// Contains on a nil Set always returns false.
func (s Set[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

// Len returns the number of elements in s.
// Len on a nil Set returns 0.
func (s Set[T]) Len() int {
	return len(s)
}

// Empty reports whether s contains no elements.
// Empty on a nil Set returns true.
func (s Set[T]) Empty() bool {
	return len(s) == 0
}

// ToSlice returns a new slice containing all elements of s in an unspecified
// order. The returned slice is never nil.
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))

	for item := range s {
		result = append(result, item)
	}

	return result
}

// String returns a human-readable representation of s in the form
// set{v1, v2, ...}. The order of elements is unspecified. A nil Set
// returns "set(nil)" and an empty Set returns "set{}".
func (s Set[T]) String() string {
	if s == nil {
		return "set(nil)"
	}

	var b strings.Builder
	b.WriteString("set{")
	first := true

	for item := range s {
		if !first {
			b.WriteString(", ")
		}

		_, _ = fmt.Fprintf(&b, "%v", item)
		first = false
	}

	b.WriteString("}")

	return b.String()
}

// Union returns a new Set containing all elements that are in s, o, or both.
// A nil operand is treated as an empty set.
func (s Set[T]) Union(o Set[T]) Set[T] {
	result := New[T]()

	for item := range s {
		result.Add(item)
	}

	for item := range o {
		result.Add(item)
	}

	return result
}

// Intersection returns a new Set containing only the elements present in
// both s and o. A nil operand is treated as an empty set.
func (s Set[T]) Intersection(o Set[T]) Set[T] {
	result := New[T]()

	for item := range s {
		if o.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Difference returns a new Set containing the elements of s that are not in o
// (set difference s∖o). A nil o is treated as an empty set (result equals s).
func (s Set[T]) Difference(o Set[T]) Set[T] {
	result := New[T]()

	for item := range s {
		if !o.Contains(item) {
			result.Add(item)
		}
	}

	return result
}

// Values returns an iterator that yields each element of s in an unspecified
// order. Values on a nil Set yields no elements.
func (s Set[T]) Values() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
