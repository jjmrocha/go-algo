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

func TestNewCountEqualsSize(t *testing.T) {
	// when
	u := New(5)
	// then
	if u.Len() != 5 {
		t.Fatalf("Expected 5 sets, got %d", u.Len())
	}
}

func TestNewEachNodeIsItsOwnRoot(t *testing.T) {
	// when
	u := New(4)
	// then
	for i := range 4 {
		root, err := u.Find(i)
		if err != nil {
			t.Fatalf("Find(%d) returned unexpected error: %v", i, err)
		}
		if root != i {
			t.Fatalf("Expected node %d to be its own root, got %d", i, root)
		}
	}
}

func TestFindOutOfRangeLow(t *testing.T) {
	// given
	u := New(3)
	// when
	_, err := u.Find(-1)
	// then
	if !errors.Is(err, ErrIndexOutOfRange) {
		t.Fatalf("Find(-1) error = %v; want ErrIndexOutOfRange", err)
	}
}

func TestFindOutOfRangeHigh(t *testing.T) {
	// given
	u := New(3)
	// when
	_, err := u.Find(3)
	// then
	if !errors.Is(err, ErrIndexOutOfRange) {
		t.Fatalf("Find(3) error = %v; want ErrIndexOutOfRange", err)
	}
}

func TestFindReturnsRoot(t *testing.T) {
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
}

func TestUnionSucceedsOnNewConnection(t *testing.T) {
	// given
	u := New(3)
	// when
	err := u.Union(0, 1)
	// then
	if err != nil {
		t.Fatalf("Union(0,1) = %v; want nil", err)
	}
}

func TestUnionSucceedsWhenAlreadyConnected(t *testing.T) {
	// given
	u := New(3)
	mustUnion(t, u, 0, 1)
	// when
	err := u.Union(0, 1)
	// then
	if err != nil {
		t.Fatalf("Union(0,1) second call = %v; want nil", err)
	}
}

func TestUnionErrorsForLowOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	err := u.Union(-1, 0)
	// then
	if err == nil {
		t.Fatalf("Union(-1,0) = nil; want ErrIndexOutOfRange")
	}
}

func TestUnionErrorsForHighOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	err := u.Union(0, 10)
	// then
	if err == nil {
		t.Fatalf("Union(0,10) = nil; want ErrIndexOutOfRange")
	}
}

func TestUnionDecreasesSetCount(t *testing.T) {
	// given
	u := New(3)
	// when
	mustUnion(t, u, 0, 1)
	// then
	if u.Len() != 2 {
		t.Fatalf("Expected 2 sets after one union, got %d", u.Len())
	}
}

func TestIsConnectedFalseForFreshNodes(t *testing.T) {
	// given
	u := New(4)
	// when
	connected, _ := u.Connected(0, 1)
	// then
	if connected {
		t.Fatalf("Expected 0 and 1 to be disconnected in a fresh QuickUnion")
	}
}

func TestIsConnectedTrueAfterUnion(t *testing.T) {
	// given
	u := New(4)
	mustUnion(t, u, 0, 1)
	// when
	connected, _ := u.Connected(0, 1)
	// then
	if !connected {
		t.Fatalf("Expected 0 and 1 to be connected after Union(0,1)")
	}
}

func TestIsConnectedTransitivity(t *testing.T) {
	// given
	u := New(4)
	mustUnion(t, u, 0, 1)
	mustUnion(t, u, 1, 2)
	// when
	connected, _ := u.Connected(0, 2)
	// then
	if !connected {
		t.Fatalf("Expected 0 and 2 to be connected via transitivity")
	}
}

func TestIsConnectedFalseForLowOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	connected, _ := u.Connected(-1, 0)
	// then
	if connected {
		t.Fatalf("Connected(-1,0) = true; want false")
	}
}

func TestIsConnectedFalseForHighOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	connected, _ := u.Connected(0, 10)
	// then
	if connected {
		t.Fatalf("Connected(0,10) = true; want false")
	}
}

func TestStringFreshForest(t *testing.T) {
	// given
	u := New(3)
	// when
	got := u.String()
	// then
	if got != "" {
		t.Fatalf("Expected empty string, got %q", got)
	}
}

func TestStringFlatTree(t *testing.T) {
	// given
	u := New(3)
	mustUnion(t, u, 0, 1)
	mustUnion(t, u, 0, 2)
	// when
	got := u.String()
	// then
	want := "0 <- 1\n0 <- 2"
	if got != want {
		t.Fatalf("Expected:\n%q\ngot:\n%q", want, got)
	}
}

func TestStringNestedTree(t *testing.T) {
	// given
	u := New(4)
	mustUnion(t, u, 0, 1)
	mustUnion(t, u, 2, 3)
	mustUnion(t, u, 0, 2)
	// when
	got := u.String()
	// then
	want := "0 <- 1\n0 <- 2\n2 <- 3"
	if got != want {
		t.Fatalf("Expected:\n%q\ngot:\n%q", want, got)
	}
}

func TestStringTwoSeparateTrees(t *testing.T) {
	// given
	u := New(6)
	mustUnion(t, u, 0, 1)
	mustUnion(t, u, 0, 2)
	mustUnion(t, u, 3, 4)
	// when
	got := u.String()
	// then
	want := "0 <- 1\n0 <- 2\n3 <- 4"
	if got != want {
		t.Fatalf("Expected:\n%q\ngot:\n%q", want, got)
	}
}
