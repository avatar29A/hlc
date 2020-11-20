package hlc

import (
	"math/bits"
	"sync"
)

// Hlc is a hybrid logical clock, that illuminated gap between the theory
// and practice of distributed systems in terms of the use of time.
//
// Based on paper: https://cse.buffalo.edu/tech-reports/2014-04.pdf
type Hlc struct {
	sync.RWMutex
	ts int64
	counter int64
	clock Clock
}

// New creates a new Hlc instance with specifics wall clock.
func New(clock Clock) *Hlc {
	return &Hlc{ts: clock.Now(), counter: 0, clock: clock}
}

// FromTimestamp takes a special NTP-compatible timestamp and returns Hlc instance
// Restriction: timestamp should contains from two part: 1. 48bit of time 2. 16 bit of counter
func FromTimestamp(ts int64) *Hlc {
	mask := uint64(0xFFFF)
	l := int64(uint64(ts) & ^mask)
	c := int64(uint64(ts) & mask)

	return &Hlc{ts: l, counter: c, clock: nil}
}

// Now returns NTP-compatible timestamp with logical time.
// Should be called for every local or outgoing event.
func (hlc *Hlc) Now() int64 {
	hlc.Lock()
	ts := hlc.ts
	hlc.ts = max(ts, hlc.clock.Now())
	if hlc.ts == ts {
		hlc.counter += 1
	} else {
		hlc.counter = 0
	}
	hlc.Unlock()

	return hlc.Timestamp()
}

// Update takes NTP-compatible timestamp with logical time and update local HLC.
// Should be called for every incoming-event.
//
// rc is a remote hlc, can be extracted from NTP-timestamp (see FromTimestamp)
func (hlc *Hlc) Update(rc *Hlc) int64 {
	hlc.Lock()
	ts := hlc.ts
	hlc.ts = max(ts, rc.ts, hlc.clock.Now())
	if hlc.ts == ts && hlc.ts == rc.ts {
		hlc.counter = max(hlc.counter, rc.counter) + 1
	} else if hlc.ts == ts {
		hlc.counter += 1
	} else if hlc.ts == rc.ts {
		hlc.counter = rc.counter + 1
	} else {
		hlc.counter = 0
	}
	hlc.Unlock()

	return hlc.Timestamp()
}

// Timestamp calculates NTP compatible timestamp. Contains physical timestamp
// (accuracy upto 48bit) and logical part (upto 16 bit)
func (hlc *Hlc) Timestamp() int64 {
	hlc.RLock()
	defer hlc.RUnlock()

	mask := bits.Reverse64(^uint64(0) >> 16 )
	ts := int64(uint64(hlc.ts) & mask)

	return ts + hlc.counter
}

func max(values ...int64) int64 {
	max := int64(0)
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max
}

