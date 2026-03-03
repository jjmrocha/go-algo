// Package sets provides a generic, unordered collection of unique elements.
//
// A [Set] is backed by a map and provides O(1) amortised [Set.Add],
// [Set.Remove], and [Set.Contains] operations. Sets are not safe for
// concurrent use without external synchronisation.
package sets

// Set is an unordered collection of unique elements of type T.
//
// The zero value of Set is nil. A nil Set is safe for read operations
// ([Set.Contains], [Set.Len], [Set.Remove], [Set.ToSlice]), but calling
// [Set.Add] on a nil Set panics. Use [New] or [FromSlice] to obtain a
// ready-to-use Set.
//
// Because Set is a map type, assigning one Set variable to another creates
// an alias: both variables refer to the same underlying collection.
type Set[T comparable] map[T]struct{}

// New returns an empty, initialized Set ready for use.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// FromSlice creates a Set containing the elements of items, silently
// discarding any duplicates.
func FromSlice[T comparable](items []T) Set[T] {
	s := make(Set[T], len(items))
	s.Add(items...)
	return s
}

// Add inserts items into s. Adding an element that is already present is a
// no-op. Calling Add on a nil Set panics; use [New] or [FromSlice] to
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

// ToSlice returns a new slice containing all elements of s in an unspecified
// order. The returned slice is never nil.
func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))

	for item := range s {
		result = append(result, item)
	}

	return result
}
