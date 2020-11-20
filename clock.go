package hlc

import (
	"math/bits"
	"time"
)

type Clock interface {
	Now() int64
}

// NTPClock use only 48bit to storing time. Preferable to use with HLC implementation.
type NTPClock struct {
	c Clock
}

// Now returns rounded upto 48bit Nano-time
func (ntp *NTPClock) Now() int64 {
	if ntp.c == nil {
		ntp.c = &NanoClock{}
	}

	mask := bits.Reverse64(^uint64(0) >> 16 )
	return int64(uint64(ntp.c.Now()) & mask)
}

// Physical clock with nano-accuracy
type NanoClock struct {
}

// Now returns timestamp with nano-seconds
func (*NanoClock) Now() int64 {
	return time.Now().UnixNano()
}

// Physical clock with second-accuracy
type SecondClock struct {

}

// Now returns timestamp with seconds
func (*SecondClock) Now() int64 {
	return time.Now().Unix()
}

// FakeClock implementation for testing. returns always determinated time.
type FakeClock struct{
	ts int64
}

// Now return determinated clock's time
func (c *FakeClock) Now() int64 {
	return c.ts
}

// Tick update clock's time
func (c *FakeClock) Tick() int64 {
	c.ts += 1
	return c.ts
}
