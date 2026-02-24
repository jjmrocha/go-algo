package union

import "testing"

func TestNewCountEqualsSize(t *testing.T) {
	// when
	u := New(5)
	// then
	if u.Count() != 5 {
		t.Fatalf("Expected 5 sets, got %d", u.Count())
	}
}

func TestNewEachNodeIsItsOwnRoot(t *testing.T) {
	// when
	u := New(4)
	// then
	for i := range 4 {
		if root := u.Find(i); root != i {
			t.Fatalf("Expected node %d to be its own root, got %d", i, root)
		}
	}
}

func TestFindOutOfRangeLow(t *testing.T) {
	// given
	u := New(3)
	// when
	got := u.Find(-1)
	// then
	if got != IndexOutOfRange {
		t.Fatalf("Find(-1) = %d; want IndexOutOfRange", got)
	}
}

func TestFindOutOfRangeHigh(t *testing.T) {
	// given
	u := New(3)
	// when
	got := u.Find(3)
	// then
	if got != IndexOutOfRange {
		t.Fatalf("Find(3) = %d; want IndexOutOfRange", got)
	}
}

func TestFindReturnsRoot(t *testing.T) {
	// given
	u := New(3)
	u.Union(0, 1)
	// when
	root0 := u.Find(0)
	root1 := u.Find(1)
	// then
	if root0 != root1 {
		t.Fatalf("Expected same root after Union(0,1), got %d and %d", root0, root1)
	}
}

func TestUnionReturnsTrueOnNewConnection(t *testing.T) {
	// given
	u := New(3)
	// when
	ok := u.Union(0, 1)
	// then
	if !ok {
		t.Fatalf("Union(0,1) = false; want true")
	}
}

func TestUnionReturnsFalseWhenAlreadyConnected(t *testing.T) {
	// given
	u := New(3)
	u.Union(0, 1)
	// when
	ok := u.Union(0, 1)
	// then
	if ok {
		t.Fatalf("Union(0,1) second call = true; want false")
	}
}

func TestUnionReturnsFalseForLowOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	ok := u.Union(-1, 0)
	// then
	if ok {
		t.Fatalf("Union(-1,0) = true; want false")
	}
}

func TestUnionReturnsFalseForHighOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	ok := u.Union(0, 10)
	// then
	if ok {
		t.Fatalf("Union(0,10) = true; want false")
	}
}

func TestUnionDecreasesSetCount(t *testing.T) {
	// given
	u := New(3)
	// when
	u.Union(0, 1)
	// then
	if u.Count() != 2 {
		t.Fatalf("Expected 2 sets after one union, got %d", u.Count())
	}
}

func TestIsConnectedFalseForFreshNodes(t *testing.T) {
	// given
	u := New(4)
	// when
	connected := u.IsConnected(0, 1)
	// then
	if connected {
		t.Fatalf("Expected 0 and 1 to be disconnected in a fresh QuickUnion")
	}
}

func TestIsConnectedTrueAfterUnion(t *testing.T) {
	// given
	u := New(4)
	u.Union(0, 1)
	// when
	connected := u.IsConnected(0, 1)
	// then
	if !connected {
		t.Fatalf("Expected 0 and 1 to be connected after Union(0,1)")
	}
}

func TestIsConnectedTransitivity(t *testing.T) {
	// given
	u := New(4)
	u.Union(0, 1)
	u.Union(1, 2)
	// when
	connected := u.IsConnected(0, 2)
	// then
	if !connected {
		t.Fatalf("Expected 0 and 2 to be connected via transitivity")
	}
}

func TestIsConnectedFalseForLowOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	connected := u.IsConnected(-1, 0)
	// then
	if connected {
		t.Fatalf("IsConnected(-1,0) = true; want false")
	}
}

func TestIsConnectedFalseForHighOutOfRange(t *testing.T) {
	// given
	u := New(3)
	// when
	connected := u.IsConnected(0, 10)
	// then
	if connected {
		t.Fatalf("IsConnected(0,10) = true; want false")
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
	u.Union(0, 1)
	u.Union(0, 2)
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
	u.Union(0, 1)
	u.Union(2, 3)
	u.Union(0, 2)
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
	u.Union(0, 1)
	u.Union(0, 2)
	u.Union(3, 4)
	// when
	got := u.String()
	// then
	want := "0 <- 1\n0 <- 2\n3 <- 4"
	if got != want {
		t.Fatalf("Expected:\n%q\ngot:\n%q", want, got)
	}
}
