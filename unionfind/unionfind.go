package unionfind

// UnionFind is the interface for union-find algorithms
type UnionFind interface {
	// Union connects two elements
	Union(p, q int)
	// Find returns the component identifier for the element
	Find(p int) int
	// Connected returns true if two elements are in the same component
	Connected(p, q int) bool
	// Count returns the number of components
	Count() int
}
