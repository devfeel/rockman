package standard

import (
	"sync/atomic"
	"time"
)

// StandardCounter is the standard implementation of a Counter
type StandardCounter struct {
	count     int64
	startTime time.Time
}

func NewStandardCounter() *StandardCounter {
	return &StandardCounter{startTime: time.Now()}
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

// Dec decrements the counter by 1.
func (c *StandardCounter) Dec() {
	c.Add(-1)
}

// Inc increments the counter by 1.
func (c *StandardCounter) Inc() {
	c.Add(1)
}

// Add increments the counter by the given value.
func (c *StandardCounter) Add(value int64) {
	atomic.AddInt64(&c.count, value)
}
