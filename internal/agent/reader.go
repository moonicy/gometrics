package agent

import (
	"github.com/shirou/gopsutil/v4/cpu"
	gopsutil "github.com/shirou/gopsutil/v4/mem"
	"log"
	"math/rand"
	"runtime"
	"strconv"
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
	mem.Gauge[Alloc] = float64(mr.rtm.Alloc)
	mem.Gauge[BuckHashSys] = float64(mr.rtm.BuckHashSys)
	mem.Gauge[Frees] = float64(mr.rtm.Frees)
	mem.Gauge[GCCPUFraction] = mr.rtm.GCCPUFraction
	mem.Gauge[GCSys] = float64(mr.rtm.GCSys)
	mem.Gauge[HeapAlloc] = float64(mr.rtm.HeapAlloc)
	mem.Gauge[HeapIdle] = float64(mr.rtm.HeapIdle)
	mem.Gauge[HeapInuse] = float64(mr.rtm.HeapInuse)
	mem.Gauge[HeapObjects] = float64(mr.rtm.HeapObjects)
	mem.Gauge[HeapReleased] = float64(mr.rtm.HeapReleased)
	mem.Gauge[HeapSys] = float64(mr.rtm.HeapSys)
	mem.Gauge[LastGC] = float64(mr.rtm.LastGC)
	mem.Gauge[Lookups] = float64(mr.rtm.Lookups)
	mem.Gauge[MCacheInuse] = float64(mr.rtm.MCacheInuse)
	mem.Gauge[MCacheSys] = float64(mr.rtm.MCacheSys)
	mem.Gauge[MSpanInuse] = float64(mr.rtm.MSpanInuse)
	mem.Gauge[MSpanSys] = float64(mr.rtm.MSpanSys)
	mem.Gauge[Mallocs] = float64(mr.rtm.Mallocs)
	mem.Gauge[NextGC] = float64(mr.rtm.NextGC)
	mem.Gauge[NumForcedGC] = float64(mr.rtm.NumForcedGC)
	mem.Gauge[NumGC] = float64(mr.rtm.NumGC)
	mem.Gauge[OtherSys] = float64(mr.rtm.OtherSys)
	mem.Gauge[PauseTotalNs] = float64(mr.rtm.PauseTotalNs)
	mem.Gauge[StackInuse] = float64(mr.rtm.StackInuse)
	mem.Gauge[StackSys] = float64(mr.rtm.StackSys)
	mem.Gauge[Sys] = float64(mr.rtm.Sys)
	mem.Gauge[TotalAlloc] = float64(mr.rtm.TotalAlloc)

	mem.Gauge[RandomValue] = rand.Float64()

	mem.Counter[PollCount]++

	mem.Gauge[TotalMemory] = float64(v.Total)
	mem.Gauge[FreeMemory] = float64(v.Free)
	cpuAll, err := cpu.Percent(0, true)
	if err != nil {
		log.Println("Error getting CPU stats: ", err)
		return
	}
	for i, cpuVal := range cpuAll {
		mem.Gauge[CPUutilization+strconv.Itoa(i+1)] = cpuVal
	}
}
