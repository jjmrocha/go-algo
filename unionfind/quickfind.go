package unionfind

// QuickFind is a simple union-find implementation with quick find operation
type QuickFind struct {
	id    []int
	count int
}

// NewQuickFind creates a new QuickFind with n elements
func NewQuickFind(n int) *QuickFind {
	id := make([]int, n)
	for i := 0; i < n; i++ {
		id[i] = i
	}
	return &QuickFind{
		id:    id,
		count: n,
	}
}

// Find returns the component identifier for the element
func (qf *QuickFind) Find(p int) int {
	return qf.id[p]
}

// Union connects two elements
func (qf *QuickFind) Union(p, q int) {
	pID := qf.Find(p)
	qID := qf.Find(q)

	if pID == qID {
		return
	}

	// Change all entries with id[p] to id[q]
	for i := 0; i < len(qf.id); i++ {
		if qf.id[i] == pID {
			qf.id[i] = qID
		}
	}
	qf.count--
}

// Connected returns true if two elements are in the same component
func (qf *QuickFind) Connected(p, q int) bool {
	return qf.Find(p) == qf.Find(q)
}

// Count returns the number of components
func (qf *QuickFind) Count() int {
	return qf.count
}
