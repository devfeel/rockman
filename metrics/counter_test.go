package metrics

import (
	"testing"
)

func TestStandardCounter_Inc(t *testing.T) {
	var wantCount int64 = 1
	counter := NewCounter()
	counter.Inc(1)
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}

func TestStandardCounter_Dec(t *testing.T) {
	var wantCount int64 = -1
	counter := NewCounter()
	counter.Dec(1)
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}

func TestStandardCounter_Clear(t *testing.T) {
	var wantCount int64 = 0
	counter := NewCounter()
	counter.Inc(1)
	counter.Inc(1)
	counter.Inc(1)
	counter.Clear()
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}
