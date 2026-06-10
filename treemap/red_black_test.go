package treemap

import (
	"math/rand"
	stdsort "sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func intCmp(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func TestNew(t *testing.T) {
	// when
	tr := New[int, int](intCmp)
	// then
	assert.Equal(t, 0, tr.Len())
	assert.True(t, tr.Empty())
}

func TestGet(t *testing.T) {
	t.Run("empty tree returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		result, ok := tr.Get(1)
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})

	t.Run("absent key returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		result, ok := tr.Get(2)
		// then
		assert.False(t, ok)
		assert.Equal(t, "", result)
	})

	t.Run("present key returns value and true", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		result, ok := tr.Get(1)
		// then
		assert.True(t, ok)
		assert.Equal(t, "a", result)
	})

	t.Run("returns updated value after overwrite", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(1, "b")
		// when
		result, ok := tr.Get(1)
		// then
		assert.True(t, ok)
		assert.Equal(t, "b", result)
	})
}

func TestContains(t *testing.T) {
	t.Run("false for empty tree", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		result := tr.Contains(1)
		// then
		assert.False(t, result)
	})

	t.Run("false for absent key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.Contains(2)
		// then
		assert.False(t, result)
	})

	t.Run("true for present key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.Contains(1)
		// then
		assert.True(t, result)
	})
}

func TestLen(t *testing.T) {
	t.Run("zero when empty", func(t *testing.T) {
		// given / when
		tr := New[int, int](intCmp)
		// then
		assert.Equal(t, 0, tr.Len())
	})

	t.Run("counts unique keys", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(2, 2)
		tr.Put(3, 3)
		// when
		result := tr.Len()
		// then
		assert.Equal(t, 3, result)
	})

	t.Run("unchanged when overwriting existing key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		tr.Put(1, 99)
		// then
		assert.Equal(t, 1, tr.Len())
	})
}

func TestEmpty(t *testing.T) {
	t.Run("true when new", func(t *testing.T) {
		// given / when
		tr := New[int, int](intCmp)
		// then
		assert.True(t, tr.Empty())
	})

	t.Run("false after put", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		tr.Put(1, 1)
		// then
		assert.False(t, tr.Empty())
	})
}

func TestPut(t *testing.T) {
	t.Run("inserts new key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		tr.Put(42, 99)
		// then
		v, ok := tr.Get(42)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})

	t.Run("overwrites value for existing key", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		tr.Put(1, 99)
		// then
		v, ok := tr.Get(1)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})
}

func TestMin(t *testing.T) {
	t.Run("empty tree returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		key, ok := tr.Min()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, key)
	})

	t.Run("single element returns that key and true", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(5, "a")
		// when
		key, ok := tr.Min()
		// then
		assert.True(t, ok)
		assert.Equal(t, 5, key)
	})

	t.Run("returns smallest key after ascending insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(2, "b")
		tr.Put(3, "c")
		// when
		key, ok := tr.Min()
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, key)
	})

	t.Run("returns smallest key after descending insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(3, "c")
		tr.Put(2, "b")
		tr.Put(1, "a")
		// when
		key, ok := tr.Min()
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, key)
	})

	t.Run("returns smallest key after random insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, "v")
		}
		// when
		key, ok := tr.Min()
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, key)
	})
}

func TestMax(t *testing.T) {
	t.Run("empty tree returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		key, ok := tr.Max()
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, key)
	})

	t.Run("single element returns that key and true", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(5, "a")
		// when
		key, ok := tr.Max()
		// then
		assert.True(t, ok)
		assert.Equal(t, 5, key)
	})

	t.Run("returns largest key after ascending insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(2, "b")
		tr.Put(3, "c")
		// when
		key, ok := tr.Max()
		// then
		assert.True(t, ok)
		assert.Equal(t, 3, key)
	})

	t.Run("returns largest key after descending insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(3, "c")
		tr.Put(2, "b")
		tr.Put(1, "a")
		// when
		key, ok := tr.Max()
		// then
		assert.True(t, ok)
		assert.Equal(t, 3, key)
	})

	t.Run("returns largest key after random insertions", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, "v")
		}
		// when
		key, ok := tr.Max()
		// then
		assert.True(t, ok)
		assert.Equal(t, 7, key)
	})
}

func TestRank(t *testing.T) {
	t.Run("empty tree returns 0", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when / then
		assert.Equal(t, 0, tr.Rank(5))
	})

	t.Run("minimum key returns 0", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{3, 1, 5} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 0, tr.Rank(1))
	})

	t.Run("maximum key returns size minus 1", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{3, 1, 5} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 2, tr.Rank(5))
	})

	t.Run("key less than all elements returns 0", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{3, 5, 7} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 0, tr.Rank(1))
	})

	t.Run("key greater than all elements returns size", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{3, 1, 5} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 3, tr.Rank(10))
	})

	t.Run("middle key returns correct rank", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{1, 2, 3, 4, 5} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 2, tr.Rank(3))
	})

	t.Run("absent key between elements returns number of keys less than it", func(t *testing.T) {
		// given — keys {1, 3, 5}, querying 2
		tr := New[int, string](intCmp)
		for _, k := range []int{1, 3, 5} {
			tr.Put(k, "v")
		}
		// when / then
		assert.Equal(t, 1, tr.Rank(2))
	})
}

func TestSelect(t *testing.T) {
	t.Run("empty tree returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		key, ok := tr.Select(0)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, key)
	})

	t.Run("negative rank returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		key, ok := tr.Select(-1)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, key)
	})

	t.Run("rank equal to size returns zero and false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		// when
		key, ok := tr.Select(1)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, key)
	})

	t.Run("rank 0 returns minimum key", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, "v")
		}
		// when
		key, ok := tr.Select(0)
		// then
		assert.True(t, ok)
		assert.Equal(t, 1, key)
	})

	t.Run("rank size-1 returns maximum key", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, "v")
		}
		// when
		key, ok := tr.Select(6)
		// then
		assert.True(t, ok)
		assert.Equal(t, 7, key)
	})

	t.Run("middle rank returns correct key", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, "v")
		}
		// when
		key, ok := tr.Select(3)
		// then
		assert.True(t, ok)
		assert.Equal(t, 4, key)
	})
}

func TestRankAndSelectAreInverses(t *testing.T) {
	// given
	tr := New[int, string](intCmp)
	keys := []int{5, 3, 7, 1, 4, 6, 2}
	for _, k := range keys {
		tr.Put(k, "v")
	}

	t.Run("Select(Rank(key)) == key for all keys in tree", func(t *testing.T) {
		for _, k := range keys {
			rank := tr.Rank(k)
			selected, ok := tr.Select(rank)
			assert.True(t, ok)
			assert.Equal(t, k, selected)
		}
	})

	t.Run("Rank(Select(r)) == r for all valid ranks", func(t *testing.T) {
		for r := range tr.Len() {
			key, ok := tr.Select(r)
			assert.True(t, ok)
			assert.Equal(t, r, tr.Rank(key))
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("empty tree returns false", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		// when
		ok := tr.Delete(1)
		// then
		assert.False(t, ok)
		assert.Equal(t, 0, tr.Len())
	})

	t.Run("absent key returns false and leaves tree unchanged", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(3, "c")
		// when
		ok := tr.Delete(2)
		// then
		assert.False(t, ok)
		assert.Equal(t, 2, tr.Len())
		assert.True(t, tr.Contains(1))
		assert.True(t, tr.Contains(3))
	})

	t.Run("present key returns true and is removed", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(2, "b")
		// when
		ok := tr.Delete(1)
		// then
		assert.True(t, ok)
		assert.False(t, tr.Contains(1))
		assert.Equal(t, 1, tr.Len())
	})

	t.Run("deleting the only element empties the tree", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(5, "a")
		// when
		ok := tr.Delete(5)
		// then
		assert.True(t, ok)
		assert.True(t, tr.Empty())
		assert.Empty(t, tr.ToList())
	})

	t.Run("deleting the same key twice returns false the second time", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(2, "b")
		// when
		first := tr.Delete(1)
		second := tr.Delete(1)
		// then
		assert.True(t, first)
		assert.False(t, second)
		assert.Equal(t, 1, tr.Len())
	})

	t.Run("deleting min then querying order", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		ok := tr.Delete(1)
		// then
		assert.True(t, ok)
		assert.Equal(t, []int{2, 3, 4, 5, 6, 7}, tr.ToList())
	})

	t.Run("deleting max then querying order", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		ok := tr.Delete(7)
		// then
		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, tr.ToList())
	})

	t.Run("deleting an internal node keeps remaining keys ordered", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		ok := tr.Delete(5)
		// then
		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3, 4, 6, 7}, tr.ToList())
		_, present := tr.Get(5)
		assert.False(t, present)
	})

	t.Run("deleting every key in sequence empties the tree and preserves invariants", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		keys := []int{5, 3, 7, 1, 4, 6, 2, 9, 8, 0}
		for _, k := range keys {
			tr.Put(k, k)
		}
		// when / then — delete in a different order, checking invariants after each delete
		for i, k := range []int{4, 0, 9, 2, 6, 5, 8, 1, 3, 7} {
			ok := tr.Delete(k)
			assert.Truef(t, ok, "expected key %d to be deleted", k)
			assert.Equalf(t, len(keys)-i-1, tr.Len(), "size after deleting %d", k)
			assertLLRBInvariants(t, tr)
		}
		assert.True(t, tr.Empty())
	})

	t.Run("remaining keys are retrievable after a delete", func(t *testing.T) {
		// given
		tr := New[int, string](intCmp)
		tr.Put(1, "a")
		tr.Put(2, "b")
		tr.Put(3, "c")
		// when
		tr.Delete(2)
		// then
		a, okA := tr.Get(1)
		c, okC := tr.Get(3)
		assert.True(t, okA)
		assert.Equal(t, "a", a)
		assert.True(t, okC)
		assert.Equal(t, "c", c)
	})

	t.Run("rank and select stay consistent after deletes", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		tr.Delete(4)
		// then — remaining sorted keys are {1,2,3,5,6,7}
		remaining := []int{1, 2, 3, 5, 6, 7}
		for r, want := range remaining {
			got, ok := tr.Select(r)
			assert.True(t, ok)
			assert.Equal(t, want, got)
			assert.Equal(t, r, tr.Rank(want))
		}
	})
}

// assertLLRBInvariants verifies the structural invariants of a Left-Leaning
// Red-Black BST: keys are in symmetric (BST) order, no right-leaning red links,
// no two consecutive left-leaning red links, and the tree is perfectly black
// balanced (every root-to-nil path has the same number of black links).
func assertLLRBInvariants[K, V any](t *testing.T, tr *Map[K, V]) {
	t.Helper()

	// BST order: in-order keys must be strictly ascending.
	prev := (*K)(nil)
	var inorder func(n *node[K, V])
	inorder = func(n *node[K, V]) {
		if n == nil {
			return
		}
		inorder(n.left)
		if prev != nil {
			assert.Negativef(t, tr.cmp(*prev, n.key), "BST order violated around key %v", n.key)
		}
		k := n.key
		prev = &k
		inorder(n.right)
	}
	inorder(tr.root)

	// No right-leaning red links and no two consecutive left red links.
	var checkRedLinks func(n *node[K, V])
	checkRedLinks = func(n *node[K, V]) {
		if n == nil {
			return
		}
		assert.Falsef(t, isRed(n.right), "right-leaning red link at key %v", n.key)
		if isRed(n) {
			assert.Falsef(t, isRed(n.left), "two consecutive left red links at key %v", n.key)
		}
		checkRedLinks(n.left)
		checkRedLinks(n.right)
	}
	checkRedLinks(tr.root)

	// Perfect black balance: count black links on the leftmost path, then assert
	// every root-to-nil path has the same count.
	expectedBlack := 0
	for n := tr.root; n != nil; n = n.left {
		if !isRed(n) {
			expectedBlack++
		}
	}
	var checkBalance func(n *node[K, V], black int)
	checkBalance = func(n *node[K, V], black int) {
		if n == nil {
			assert.Equal(t, expectedBlack, black, "black height imbalance")
			return
		}
		if !isRed(n) {
			black++
		}
		checkBalance(n.left, black)
		checkBalance(n.right, black)
	}
	checkBalance(tr.root, 0)
}

func TestDeleteStress(t *testing.T) {
	// Randomized property test: interleave random Puts and Deletes against a
	// reference map, asserting the tree stays a valid LLRB and agrees with the
	// reference on contents and ordering after every mutation.
	rng := rand.New(rand.NewSource(42))
	const keySpace = 50

	tr := New[int, int](intCmp)
	ref := make(map[int]int)

	for range 5000 {
		key := rng.Intn(keySpace)

		if rng.Intn(2) == 0 {
			value := rng.Int()
			tr.Put(key, value)
			ref[key] = value
		} else {
			_, existed := ref[key]
			ok := tr.Delete(key)
			assert.Equalf(t, existed, ok, "Delete(%d) return value", key)
			delete(ref, key)
		}

		assert.Equal(t, len(ref), tr.Len())
		assertLLRBInvariants(t, tr)

		if t.Failed() {
			t.Fatalf("invariant or size check failed after operating on key %d", key)
		}
	}

	// Final contents must match the reference exactly, in sorted order.
	wantKeys := make([]int, 0, len(ref))
	for k := range ref {
		wantKeys = append(wantKeys, k)
	}
	stdsort.Ints(wantKeys)

	wantValues := make([]int, 0, len(ref))
	for _, k := range wantKeys {
		wantValues = append(wantValues, ref[k])
	}
	assert.Equal(t, wantValues, tr.ToList())

	for k, want := range ref {
		got, ok := tr.Get(k)
		assert.Truef(t, ok, "expected key %d present", k)
		assert.Equalf(t, want, got, "value for key %d", k)
	}
}

func TestToList(t *testing.T) {
	t.Run("empty tree returns empty slice", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		// when
		result := tr.ToList()
		// then
		assert.Empty(t, result)
	})

	t.Run("single element", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1}, result)
	})

	t.Run("ascending insertion returns sorted values", func(t *testing.T) {
		// given — forces left-rotation path; exposes nil-dereference bug in naive traversal
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(2, 2)
		tr.Put(3, 3)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("descending insertion returns sorted values", func(t *testing.T) {
		// given — builds left-heavy tree before LLRBT rotations balance it
		tr := New[int, int](intCmp)
		tr.Put(3, 3)
		tr.Put(2, 2)
		tr.Put(1, 1)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("random insertion order returns sorted values", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		for _, k := range []int{5, 3, 7, 1, 4, 6, 2} {
			tr.Put(k, k)
		}
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, result)
	})

	t.Run("overwritten key appears once with latest value", func(t *testing.T) {
		// given
		tr := New[int, int](intCmp)
		tr.Put(1, 1)
		tr.Put(1, 99)
		// when
		result := tr.ToList()
		// then
		assert.Equal(t, []int{99}, result)
	})
}
