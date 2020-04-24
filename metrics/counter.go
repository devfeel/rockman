package metrics

import (
	"sync/atomic"
	"time"
)

// Counter incremented and decremented base on int64 value.
type Counter interface {
	StartTime() time.Time
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

// NewCounter constructs a new StandardCounter.
func NewCounter() Counter {
	return &StandardCounter{startTime: time.Now()}
}

// StandardCounter is the standard implementation of a Counter
type StandardCounter struct {
	count     int64
	startTime time.Time
}

func (c *StandardCounter) StartTime() time.Time {
	return c.startTime
}

// Clear sets the counter to zero.
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
