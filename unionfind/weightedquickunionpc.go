package unionfind

// WeightedQuickUnionPC is a union-find implementation with weighted quick union and path compression
type WeightedQuickUnionPC struct {
	id    []int
	sz    []int
	count int
}

// NewWeightedQuickUnionPC creates a new WeightedQuickUnionPC with n elements
func NewWeightedQuickUnionPC(n int) *WeightedQuickUnionPC {
	id := make([]int, n)
	sz := make([]int, n)
	for i := 0; i < n; i++ {
		id[i] = i
		sz[i] = 1
	}
	return &WeightedQuickUnionPC{
		id:    id,
		sz:    sz,
		count: n,
	}
}

// Find returns the component identifier for the element with path compression
func (wqupc *WeightedQuickUnionPC) Find(p int) int {
	root := p
	for root != wqupc.id[root] {
		root = wqupc.id[root]
	}
	// Path compression: make every node point to the root
	for p != root {
		next := wqupc.id[p]
		wqupc.id[p] = root
		p = next
	}
	return root
}

// Union connects two elements
func (wqupc *WeightedQuickUnionPC) Union(p, q int) {
	rootP := wqupc.Find(p)
	rootQ := wqupc.Find(q)

	if rootP == rootQ {
		return
	}

	// Link smaller tree below larger tree
	if wqupc.sz[rootP] < wqupc.sz[rootQ] {
		wqupc.id[rootP] = rootQ
		wqupc.sz[rootQ] += wqupc.sz[rootP]
	} else {
		wqupc.id[rootQ] = rootP
		wqupc.sz[rootP] += wqupc.sz[rootQ]
	}
	wqupc.count--
}

// Connected returns true if two elements are in the same component
func (wqupc *WeightedQuickUnionPC) Connected(p, q int) bool {
	return wqupc.Find(p) == wqupc.Find(q)
}

// Count returns the number of components
func (wqupc *WeightedQuickUnionPC) Count() int {
	return wqupc.count
}
