package union

import (
	"testing"
)

func TestNew_count(t *testing.T) {
	u := New(5)
	if u.Count() != 5 {
		t.Errorf("expected 5 sets, got %d", u.Count())
	}
}

func TestNew_eachNodeIsItsOwnRoot(t *testing.T) {
	u := New(4)
	for i := 0; i < 4; i++ {
		if root := u.Find(i); root != i {
			t.Errorf("Find(%d) = %d; want %d", i, root, i)
		}
	}
}

func TestFind_outOfRange(t *testing.T) {
	u := New(3)
	if got := u.Find(-1); got != IndexOutOfRange {
		t.Errorf("Find(-1) = %d; want IndexOutOfRange", got)
	}
	if got := u.Find(3); got != IndexOutOfRange {
		t.Errorf("Find(3) = %d; want IndexOutOfRange", got)
	}
}

func TestFind_returnsRoot(t *testing.T) {
	u := New(3)
	u.Union(0, 1)
	if u.Find(0) != u.Find(1) {
		t.Errorf("Find(0) != Find(1) after Union(0,1)")
	}
}

func TestUnion_returnsTrueOnNewConnection(t *testing.T) {
	u := New(3)
	if ok := u.Union(0, 1); !ok {
		t.Error("Union(0,1) = false; want true")
	}
}

func TestUnion_returnsFalseWhenAlreadyConnected(t *testing.T) {
	u := New(3)
	u.Union(0, 1)
	if ok := u.Union(0, 1); ok {
		t.Error("Union(0,1) second call = true; want false")
	}
}

func TestUnion_returnsFalseForOutOfRange(t *testing.T) {
	u := New(3)
	if ok := u.Union(-1, 0); ok {
		t.Error("Union(-1,0) = true; want false")
	}
	if ok := u.Union(0, 10); ok {
		t.Error("Union(0,10) = true; want false")
	}
}

func TestUnion_decreasesSetCount(t *testing.T) {
	u := New(4)
	u.Union(0, 1)
	if u.Count() != 3 {
		t.Errorf("Count after one union = %d; want 3", u.Count())
	}
	u.Union(2, 3)
	if u.Count() != 2 {
		t.Errorf("Count after two unions = %d; want 2", u.Count())
	}
	u.Union(0, 2)
	if u.Count() != 1 {
		t.Errorf("Count after three unions = %d; want 1", u.Count())
	}
}

func TestIsConnected_falseForFreshNodes(t *testing.T) {
	u := New(4)
	if u.IsConnected(0, 1) {
		t.Error("expected 0 and 1 to be disconnected in a fresh UnionFind")
	}
}

func TestIsConnected_trueAfterUnion(t *testing.T) {
	u := New(4)
	u.Union(0, 1)
	if !u.IsConnected(0, 1) {
		t.Error("expected 0 and 1 to be connected after Union(0,1)")
	}
}

func TestIsConnected_transitivity(t *testing.T) {
	u := New(4)
	u.Union(0, 1)
	u.Union(1, 2)
	if !u.IsConnected(0, 2) {
		t.Error("expected 0 and 2 to be connected via transitivity")
	}
}

func TestIsConnected_outOfRange(t *testing.T) {
	u := New(3)
	if u.IsConnected(-1, 0) {
		t.Error("IsConnected(-1,0) = true; want false")
	}
	if u.IsConnected(0, 10) {
		t.Error("IsConnected(0,10) = true; want false")
	}
}

func TestString_freshForest(t *testing.T) {
	u := New(3)
	if got := u.String(); got != "" {
		t.Errorf("got %q; want empty string", got)
	}
}

func TestString_flatTree(t *testing.T) {
	u := New(3)
	u.Union(0, 1)
	u.Union(0, 2)
	want := "0 <- 1 \n0 <- 2 \n"
	if got := u.String(); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestString_nestedTree(t *testing.T) {
	u := New(4)
	u.Union(0, 1)
	u.Union(2, 3)
	u.Union(0, 2)
	want := "0 <- 1 \n0 <- 2 \n2 <- 3 \n"
	if got := u.String(); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}

func TestString_twoSeparateTrees(t *testing.T) {
	u := New(6)
	u.Union(0, 1)
	u.Union(0, 2)
	u.Union(3, 4)
	want := "0 <- 1 \n0 <- 2 \n3 <- 4 \n"
	if got := u.String(); got != want {
		t.Errorf("got:\n%q\nwant:\n%q", got, want)
	}
}
