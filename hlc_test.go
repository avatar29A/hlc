package hlc

import (
	"testing"
)

func TestFromTimestamp(t *testing.T) {
	ts := int64(1605806872706744321)
	hlc := FromTimestamp(ts)

	l :=int64(1605806872706744320)
	c :=int64(1)

	if l != hlc.ts || c != hlc.counter {
		t.Errorf("unexpected result. %d != %d and %d != %d", l, hlc.ts, c, hlc.counter)
	}
}

func TestHlc_PaperExample(t *testing.T) {
	// Implementation of example from paper (3.3 HLC Algorithm, Figure 6.)
	// (see also a scheme testdata/example.png)

	// Initialise four nodes j = [0, 1, 2, 3]
	// 0 node with Initial state (10, 10, 0), others (0,0,0)

	// Initialise a shared physical clock between  1,2,3 nodes.
	pt := &FakeClock{0}

	//
	// TICK: 0

	hlc0 := New(&FakeClock{10}) // node 0 (pt: 10, l: 10, c: 0)
	hlc1 := New(pt)             // node 1 (0, 0, 0)
	hlc2 := New(pt)             // node 2 (0, 0, 0)
	hlc3 := New(pt)             // node 3 (0, 0, 0)

	//
	// TICK: 1 (pt: 1)
	pt.Tick()

	hlc1.Update(hlc0)
	hlc2.Now()
	hlc3.Now()

	// Check Invariants for TICK 1
	// hlc1: (1,10,1)
	// hlc2: (1,1,0)
	// hlc3: (1,1,0)
	compareHlc(t, hlc1, &Hlc{ts: 10, counter: 1})
	compareHlc(t, hlc2, &Hlc{ts: 1, counter: 0})
	compareHlc(t, hlc3, &Hlc{ts: 1, counter: 0})

	//
	// TICK: 2 (pt: 2)
	pt.Tick()
	hlc1.Now()
	hlc2.Update(hlc1)
	hlc3.Now()

	// Check Invariants for TICK 2
	// hlc1: (2,10,2)
	// hlc2: (2,10,3)
	// hlc3: (2,2,0)
	compareHlc(t, hlc1, &Hlc{ts: 10, counter: 2})
	compareHlc(t, hlc2, &Hlc{ts: 10, counter: 3})
	compareHlc(t, hlc3, &Hlc{ts: 2, counter: 0})

	//
	// TICK: 3 (pt: 3)
	pt.Tick()
	hlc1.Now()
	hlc2.Now()
	hlc3.Update(hlc2)

	// Check Invariants for TICK 3
	// hlc1: (3, 10, 3) paper has an error (3, 13)
	// hlc2: (3, 10, 4)
	// hlc3: (3, 10, 5)
	compareHlc(t, hlc1, &Hlc{ts: 10, counter: 3})
	compareHlc(t, hlc2, &Hlc{ts: 10, counter: 4})
	compareHlc(t, hlc3, &Hlc{ts: 10, counter: 5})

	//
	// TICK: 4 (pt: 4)
	pt.Tick()
	hlc3.Now()
	hlc2.Now()
	hlc1.Update(hlc3)

	// Check Invariants for TICK 4
	// hlc1: (4, 10, 7)
	// hlc2: (4, 10, 5)
	// hlc3: (4, 10, 6)
	compareHlc(t, hlc1, &Hlc{ts: 10, counter: 7})
	compareHlc(t, hlc2, &Hlc{ts: 10, counter: 5})
	compareHlc(t, hlc3, &Hlc{ts: 10, counter: 6})
}

func TestHlc_LocalEvents(t *testing.T) {
	pt := &FakeClock{10}
	hlc0 := New(pt) // node 0 (pt: 10, l: 10, c: 0)

	// TICK: 11
	pt.Tick()

	hlc0.Now()
	compareHlc(t, hlc0, &Hlc{ts: 11, counter: 0})

	// TICK: 12
	pt.Tick()

	hlc0.Now()
	compareHlc(t, hlc0, &Hlc{ts: 12, counter: 0})

	// TICK: 13
	pt.Tick()

	hlc0.Now()
	compareHlc(t, hlc0, &Hlc{ts: 13, counter: 0})
}

func TestHlc_Timestamp(t *testing.T) {
	lc := New(&FakeClock{ts: 1605806872706798000})
	ts := lc.Now()

	expected := int64(1605806872706744321)

	if ts != expected {
		t.Errorf("timestamp was wrong. Expected %d, but got %d", expected, ts)
	}
}

func compareHlc(t *testing.T, expected, actual *Hlc)  {
	if actual.ts != expected.ts || actual.counter != expected.counter {
		t.Fatalf("expected %v, but got %v", expected, actual)
	}
}

func TestHlc_Max(t *testing.T) {
	cases := []struct {
		Expected int64
		Result int64
	} {
		{0, max()},
		{1, max(1)},
		{2, max(2, 1)},
		{4, max(2, 1, 3, 4, 1)},
	}

	for _, c := range cases{
		if c.Result != c.Expected {
			t.Errorf("expected %d, but got %d", c.Expected, c.Result)
		}
	}
}