package agent

import (
	"log"
	"math/rand"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	gopsutil "github.com/shirou/gopsutil/v4/mem"
)

type MetricsReader struct {
	rtm runtime.MemStats
}

func NewMetricsReader() *MetricsReader {
	return &MetricsReader{}
}

func (mr *MetricsReader) Read(mem *Report) {
	v, err := gopsutil.VirtualMemory()
	if err != nil {
		log.Println("Error getting memory info", err)
		return
	}
	runtime.ReadMemStats(&mr.rtm)
	mem.SetGauge(Alloc, float64(mr.rtm.Alloc))
	mem.SetGauge(BuckHashSys, float64(mr.rtm.BuckHashSys))
	mem.SetGauge(Frees, float64(mr.rtm.Frees))
	mem.SetGauge(GCCPUFraction, mr.rtm.GCCPUFraction)
	mem.SetGauge(GCSys, float64(mr.rtm.GCSys))
	mem.SetGauge(HeapAlloc, float64(mr.rtm.HeapAlloc))
	mem.SetGauge(HeapIdle, float64(mr.rtm.HeapIdle))
	mem.SetGauge(HeapInuse, float64(mr.rtm.HeapInuse))
	mem.SetGauge(HeapObjects, float64(mr.rtm.HeapObjects))
	mem.SetGauge(HeapReleased, float64(mr.rtm.HeapReleased))
	mem.SetGauge(HeapSys, float64(mr.rtm.HeapSys))
	mem.SetGauge(LastGC, float64(mr.rtm.LastGC))
	mem.SetGauge(Lookups, float64(mr.rtm.Lookups))
	mem.SetGauge(MCacheInuse, float64(mr.rtm.MCacheInuse))
	mem.SetGauge(MCacheSys, float64(mr.rtm.MCacheSys))
	mem.SetGauge(MSpanInuse, float64(mr.rtm.MSpanInuse))
	mem.SetGauge(MSpanSys, float64(mr.rtm.MSpanSys))
	mem.SetGauge(Mallocs, float64(mr.rtm.Mallocs))
	mem.SetGauge(NextGC, float64(mr.rtm.NextGC))
	mem.SetGauge(NumForcedGC, float64(mr.rtm.NumForcedGC))
	mem.SetGauge(NumGC, float64(mr.rtm.NumGC))
	mem.SetGauge(OtherSys, float64(mr.rtm.OtherSys))
	mem.SetGauge(PauseTotalNs, float64(mr.rtm.PauseTotalNs))
	mem.SetGauge(StackInuse, float64(mr.rtm.StackInuse))
	mem.SetGauge(StackSys, float64(mr.rtm.StackSys))
	mem.SetGauge(Sys, float64(mr.rtm.Sys))
	mem.SetGauge(TotalAlloc, float64(mr.rtm.TotalAlloc))
	mem.SetGauge(RandomValue, rand.Float64())

	mem.AddCounter(PollCount, 1)

	mem.SetGauge(TotalMemory, float64(v.Total))
	mem.SetGauge(FreeMemory, float64(v.Free))
	cpuAll, err := cpu.Percent(0, true)
	if err != nil {
		log.Println("Error getting CPU stats: ", err)
		return
	}
	for i, cpuVal := range cpuAll {
		mem.SetGauge(CPUutilization+strconv.Itoa(i+1), cpuVal)
	}
}
