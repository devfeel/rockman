package metrics

import "testing"

func TestStandardMetrics_GetDefaultCounter(t *testing.T) {
	var needVal int64 = 1
	m := NewMetrics()
	c := m.GetDefaultCounter()
	c.Inc(1)
	c2 := m.GetDefaultCounter()
	if c2.Count() == needVal {
		t.Log("success")
	} else {
		t.Error("val not need", needVal, c2.Count())
	}
}
