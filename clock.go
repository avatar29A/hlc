package hlc

import (
	"time"
)

type Clock interface {
	Now() int64
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
