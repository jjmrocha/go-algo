package unionfind

// WeightedQuickUnion is a union-find implementation with weighted quick union
type WeightedQuickUnion struct {
	id    []int
	sz    []int
	count int
}

// NewWeightedQuickUnion creates a new WeightedQuickUnion with n elements
func NewWeightedQuickUnion(n int) *WeightedQuickUnion {
	id := make([]int, n)
	sz := make([]int, n)
	for i := 0; i < n; i++ {
		id[i] = i
		sz[i] = 1
	}
	return &WeightedQuickUnion{
		id:    id,
		sz:    sz,
		count: n,
	}
}

// Find returns the component identifier for the element
func (wqu *WeightedQuickUnion) Find(p int) int {
	for p != wqu.id[p] {
		p = wqu.id[p]
	}
	return p
}

// Union connects two elements
func (wqu *WeightedQuickUnion) Union(p, q int) {
	rootP := wqu.Find(p)
	rootQ := wqu.Find(q)

	if rootP == rootQ {
		return
	}

	// Link smaller tree below larger tree
	if wqu.sz[rootP] < wqu.sz[rootQ] {
		wqu.id[rootP] = rootQ
		wqu.sz[rootQ] += wqu.sz[rootP]
	} else {
		wqu.id[rootQ] = rootP
		wqu.sz[rootP] += wqu.sz[rootQ]
	}
	wqu.count--
}

// Connected returns true if two elements are in the same component
func (wqu *WeightedQuickUnion) Connected(p, q int) bool {
	return wqu.Find(p) == wqu.Find(q)
}

// Count returns the number of components
func (wqu *WeightedQuickUnion) Count() int {
	return wqu.count
}
