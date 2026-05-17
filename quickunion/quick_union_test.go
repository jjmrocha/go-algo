package quickunion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustUnion(t *testing.T, u *QuickUnion, p, q int) {
	t.Helper()
	err := u.Union(p, q)
	if err != nil {
		t.Fatalf("setup Union(%d, %d) failed: %v", p, q, err)
	}
}

func TestNew(t *testing.T) {
	t.Run("count equals size", func(t *testing.T) {
		// when
		result := New(5)
		// then
		assert.Equal(t, 5, result.Len())
	})

	t.Run("each node is its own root", func(t *testing.T) {
		// when
		result := New(4)
		// then
		for i := range 4 {
			root, err := result.Find(i)
			assert.NoError(t, err)
			assert.Equal(t, i, root)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("out of range low", func(t *testing.T) {
		// given
		u := New(3)
		// when
		_, result := u.Find(-1)
		// then
		assert.ErrorIs(t, result, ErrIndexOutOfRange)
	})

	t.Run("out of range high", func(t *testing.T) {
		// given
		u := New(3)
		// when
		_, result := u.Find(3)
		// then
		assert.ErrorIs(t, result, ErrIndexOutOfRange)
	})

	t.Run("returns root after unionfind", func(t *testing.T) {
		// given
		u := New(3)
		mustUnion(t, u, 0, 1)
		// when
		root0, _ := u.Find(0)
		root1, _ := u.Find(1)
		// then
		assert.Equal(t, root0, root1)
	})
}

func TestUnion(t *testing.T) {
	t.Run("succeeds on new connection", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(0, 1)
		// then
		assert.NoError(t, result)
	})

	t.Run("succeeds when already connected", func(t *testing.T) {
		// given
		u := New(3)
		mustUnion(t, u, 0, 1)
		// when
		result := u.Union(0, 1)
		// then
		assert.NoError(t, result)
	})

	t.Run("errors for low out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(-1, 0)
		// then
		assert.ErrorIs(t, result, ErrIndexOutOfRange)
	})

	t.Run("errors for high out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(0, 10)
		// then
		assert.ErrorIs(t, result, ErrIndexOutOfRange)
	})

	t.Run("decreases set count", func(t *testing.T) {
		// given
		u := New(3)
		// when
		mustUnion(t, u, 0, 1)
		// then
		assert.Equal(t, 2, u.Len())
	})
}

func TestConnected(t *testing.T) {
	t.Run("false for fresh nodes", func(t *testing.T) {
		// given
		u := New(4)
		// when
		result, _ := u.Connected(0, 1)
		// then
		assert.False(t, result)
	})

	t.Run("true after unionfind", func(t *testing.T) {
		// given
		u := New(4)
		mustUnion(t, u, 0, 1)
		// when
		result, _ := u.Connected(0, 1)
		// then
		assert.True(t, result)
	})

	t.Run("transitivity", func(t *testing.T) {
		// given
		u := New(4)
		mustUnion(t, u, 0, 1)
		mustUnion(t, u, 1, 2)
		// when
		result, _ := u.Connected(0, 2)
		// then
		assert.True(t, result)
	})

	t.Run("false for low out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result, _ := u.Connected(-1, 0)
		// then
		assert.False(t, result)
	})

	t.Run("false for high out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result, _ := u.Connected(0, 10)
		// then
		assert.False(t, result)
	})
}

func TestString(t *testing.T) {
	t.Run("fresh forest", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.String()
		// then
		assert.Empty(t, result)
	})

	t.Run("flat tree", func(t *testing.T) {
		// given
		u := New(3)
		mustUnion(t, u, 0, 1)
		mustUnion(t, u, 0, 2)
		// when
		result := u.String()
		// then
		expected := "0 <- 1\n0 <- 2"
		assert.Equal(t, expected, result)
	})

	t.Run("nested tree", func(t *testing.T) {
		// given
		u := New(4)
		mustUnion(t, u, 0, 1)
		mustUnion(t, u, 2, 3)
		mustUnion(t, u, 0, 2)
		// when
		result := u.String()
		// then
		expected := "0 <- 1\n0 <- 2\n2 <- 3"
		assert.Equal(t, expected, result)
	})

	t.Run("two separate trees", func(t *testing.T) {
		// given
		u := New(6)
		mustUnion(t, u, 0, 1)
		mustUnion(t, u, 0, 2)
		mustUnion(t, u, 3, 4)
		// when
		result := u.String()
		// then
		expected := "0 <- 1\n0 <- 2\n3 <- 4"
		assert.Equal(t, expected, result)
	})
}
