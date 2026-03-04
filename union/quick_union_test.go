package union

import (
	"errors"
	"testing"
)

func mustUnion(t *testing.T, u *QuickUnion, p, q int) {
	t.Helper()
	if err := u.Union(p, q); err != nil {
		t.Fatalf("setup Union(%d, %d) failed: %v", p, q, err)
	}
}

func TestNew(t *testing.T) {
	t.Run("count equals size", func(t *testing.T) {
		// when
		result := New(5)
		// then
		if result.Len() != 5 {
			t.Fatalf("Expected 5 sets, got %d", result.Len())
		}
	})

	t.Run("each node is its own root", func(t *testing.T) {
		// when
		result := New(4)
		// then
		for i := range 4 {
			root, err := result.Find(i)
			if err != nil {
				t.Fatalf("Find(%d) returned unexpected error: %v", i, err)
			}
			if root != i {
				t.Fatalf("Expected node %d to be its own root, got %d", i, root)
			}
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
		if !errors.Is(result, ErrIndexOutOfRange) {
			t.Fatalf("Find(-1) error = %v; want ErrIndexOutOfRange", result)
		}
	})

	t.Run("out of range high", func(t *testing.T) {
		// given
		u := New(3)
		// when
		_, result := u.Find(3)
		// then
		if !errors.Is(result, ErrIndexOutOfRange) {
			t.Fatalf("Find(3) error = %v; want ErrIndexOutOfRange", result)
		}
	})

	t.Run("returns root after union", func(t *testing.T) {
		// given
		u := New(3)
		mustUnion(t, u, 0, 1)
		// when
		root0, _ := u.Find(0)
		root1, _ := u.Find(1)
		// then
		if root0 != root1 {
			t.Fatalf("Expected same root after Union(0,1), got %d and %d", root0, root1)
		}
	})
}

func TestUnion(t *testing.T) {
	t.Run("succeeds on new connection", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(0, 1)
		// then
		if result != nil {
			t.Fatalf("Union(0,1) = %v; want nil", result)
		}
	})

	t.Run("succeeds when already connected", func(t *testing.T) {
		// given
		u := New(3)
		mustUnion(t, u, 0, 1)
		// when
		result := u.Union(0, 1)
		// then
		if result != nil {
			t.Fatalf("Union(0,1) second call = %v; want nil", result)
		}
	})

	t.Run("errors for low out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(-1, 0)
		// then
		if result == nil {
			t.Fatalf("Union(-1,0) = nil; want ErrIndexOutOfRange")
		}
	})

	t.Run("errors for high out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.Union(0, 10)
		// then
		if result == nil {
			t.Fatalf("Union(0,10) = nil; want ErrIndexOutOfRange")
		}
	})

	t.Run("decreases set count", func(t *testing.T) {
		// given
		u := New(3)
		// when
		mustUnion(t, u, 0, 1)
		// then
		if u.Len() != 2 {
			t.Fatalf("Expected 2 sets after one union, got %d", u.Len())
		}
	})
}

func TestConnected(t *testing.T) {
	t.Run("false for fresh nodes", func(t *testing.T) {
		// given
		u := New(4)
		// when
		result, _ := u.Connected(0, 1)
		// then
		if result {
			t.Fatalf("Expected 0 and 1 to be disconnected in a fresh QuickUnion")
		}
	})

	t.Run("true after union", func(t *testing.T) {
		// given
		u := New(4)
		mustUnion(t, u, 0, 1)
		// when
		result, _ := u.Connected(0, 1)
		// then
		if !result {
			t.Fatalf("Expected 0 and 1 to be connected after Union(0,1)")
		}
	})

	t.Run("transitivity", func(t *testing.T) {
		// given
		u := New(4)
		mustUnion(t, u, 0, 1)
		mustUnion(t, u, 1, 2)
		// when
		result, _ := u.Connected(0, 2)
		// then
		if !result {
			t.Fatalf("Expected 0 and 2 to be connected via transitivity")
		}
	})

	t.Run("false for low out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result, _ := u.Connected(-1, 0)
		// then
		if result {
			t.Fatalf("Connected(-1,0) = true; want false")
		}
	})

	t.Run("false for high out of range", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result, _ := u.Connected(0, 10)
		// then
		if result {
			t.Fatalf("Connected(0,10) = true; want false")
		}
	})
}

func TestString(t *testing.T) {
	t.Run("fresh forest", func(t *testing.T) {
		// given
		u := New(3)
		// when
		result := u.String()
		// then
		if result != "" {
			t.Fatalf("Expected empty string, got %q", result)
		}
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
		if result != expected {
			t.Fatalf("Expected:\n%q\ngot:\n%q", expected, result)
		}
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
		if result != expected {
			t.Fatalf("Expected:\n%q\ngot:\n%q", expected, result)
		}
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
		if result != expected {
			t.Fatalf("Expected:\n%q\ngot:\n%q", expected, result)
		}
	})
}
