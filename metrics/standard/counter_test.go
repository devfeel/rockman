package standard

import (
	"testing"
)

func TestStandardCounter_Inc(t *testing.T) {
	var wantCount int64 = 1
	counter := NewStandardCounter()
	counter.Inc()
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}

func TestStandardCounter_Dec(t *testing.T) {
	var wantCount int64 = -1
	counter := NewStandardCounter()
	counter.Dec()
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}

func TestStandardCounter_Clear(t *testing.T) {
	var wantCount int64 = 0
	counter := NewStandardCounter()
	counter.Inc()
	counter.Inc()
	counter.Inc()
	counter.Clear()
	realCount := counter.Count()
	if realCount != wantCount {
		t.Error("count not correct:", wantCount, realCount)
	} else {
		t.Log("success.")
	}
}
