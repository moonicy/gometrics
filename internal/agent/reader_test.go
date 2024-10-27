package agent

import (
	"testing"
)

func TestMetricsReader_Read(t *testing.T) {
	mr := NewMetricsReader()
	mem := NewReport()
	mr.Read(mem)
	if _, exist := mem.gauge[Alloc]; !exist {
		t.Errorf("gauge wasn't filled")
	}
}
