package unionfind

// QuickUnion is a union-find implementation with quick union operation
type QuickUnion struct {
	id    []int
	count int
}

// NewQuickUnion creates a new QuickUnion with n elements
func NewQuickUnion(n int) *QuickUnion {
	id := make([]int, n)
	for i := 0; i < n; i++ {
		id[i] = i
	}
	return &QuickUnion{
		id:    id,
		count: n,
	}
}

// Find returns the component identifier for the element
func (qu *QuickUnion) Find(p int) int {
	for p != qu.id[p] {
		p = qu.id[p]
	}
	return p
}

// Union connects two elements
func (qu *QuickUnion) Union(p, q int) {
	rootP := qu.Find(p)
	rootQ := qu.Find(q)

	if rootP == rootQ {
		return
	}

	qu.id[rootP] = rootQ
	qu.count--
}

// Connected returns true if two elements are in the same component
func (qu *QuickUnion) Connected(p, q int) bool {
	return qu.Find(p) == qu.Find(q)
}

// Count returns the number of components
func (qu *QuickUnion) Count() int {
	return qu.count
}
