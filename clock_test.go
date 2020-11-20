package hlc

import (
	"testing"
)

func TestNTPClock_Now(t *testing.T) {
	clock := NTPClock{&FakeClock{ts: 1605806872706798000}}
	ts := clock.Now()
	expected := int64(1605806872706744320)

	if ts != expected {
		t.Errorf("bad rounded time. expected %d, but got %d", expected, ts)
	}
}
