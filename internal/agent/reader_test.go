package agent

import (
	"testing"
)

func TestMetricsReader_Read(t *testing.T) {
	mr := NewMetricsReader()
	mem := NewReport()
	mr.Read(mem)
	if _, exist := mem.Gauge.Load(Alloc); !exist {
		t.Errorf("gauge wasn't filled")
	}
}
