package metrics

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
)

const (
	DefaultDiskPath            = "/"
	DefaultCPUSamplingInterval = 200 * time.Millisecond
)

var (
	ErrInvalidCPUSamplingInterval = errors.New("cpu sampling interval must be greater than zero")
	ErrInvalidDiskPath            = errors.New("disk path must not be empty")
)

type CPUMetrics struct {
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryMetrics struct {
	TotalBytes       uint64  `json:"total_bytes"`
	UsedBytes        uint64  `json:"used_bytes"`
	FreeBytes        uint64  `json:"free_bytes"`
	AvailableBytes   uint64  `json:"available_bytes"`
	UsagePercent     float64 `json:"usage_percent"`
	SwapTotalBytes   uint64  `json:"swap_total_bytes"`
	SwapUsedBytes    uint64  `json:"swap_used_bytes"`
	SwapFreeBytes    uint64  `json:"swap_free_bytes"`
	SwapUsagePercent float64 `json:"swap_usage_percent"`
}

type DiskMetrics struct {
	Path         string  `json:"path"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type LoadMetrics struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

type GoRuntimeMetrics struct {
	GoroutineCount     int    `json:"goroutine_count"`
	GOMAXPROCS         int    `json:"gomaxprocs"`
	CgoCallCount       int64  `json:"cgo_call_count"`
	AllocBytes         uint64 `json:"alloc_bytes"`
	TotalAllocBytes    uint64 `json:"total_alloc_bytes"`
	SysBytes           uint64 `json:"sys_bytes"`
	HeapAllocBytes     uint64 `json:"heap_alloc_bytes"`
	HeapSysBytes       uint64 `json:"heap_sys_bytes"`
	HeapInuseBytes     uint64 `json:"heap_inuse_bytes"`
	StackInuseBytes    uint64 `json:"stack_inuse_bytes"`
	NextGCBytes        uint64 `json:"next_gc_bytes"`
	LastGCTimeUnixNano uint64 `json:"last_gc_time_unix_nano"`
	NumGC              uint32 `json:"num_gc"`
	PauseTotalNs       uint64 `json:"pause_total_ns"`
}

type MachineMetrics struct {
	CPU       CPUMetrics       `json:"cpu"`
	Memory    MemoryMetrics    `json:"memory"`
	Disk      DiskMetrics      `json:"disk"`
	Load      LoadMetrics      `json:"load"`
	GoRuntime GoRuntimeMetrics `json:"go_runtime"`
}

type CollectOptions struct {
	DiskPath            string
	CPUSamplingInterval time.Duration
}

func GetCPUMetrics() (CPUMetrics, error) {
	return GetCPUMetricsWithInterval(DefaultCPUSamplingInterval)
}

func GetCPUMetricsWithInterval(interval time.Duration) (CPUMetrics, error) {
	if interval <= 0 {
		return CPUMetrics{}, ErrInvalidCPUSamplingInterval
	}

	before, err := cpu.Get()
	if err != nil {
		return CPUMetrics{}, fmt.Errorf("collect cpu metrics: %w", err)
	}
	time.Sleep(interval)
	after, err := cpu.Get()
	if err != nil {
		return CPUMetrics{}, fmt.Errorf("collect cpu metrics: %w", err)
	}

	return CPUMetrics{
		UsagePercent: cpuUsagePercent(before.Total, before.Idle, after.Total, after.Idle),
	}, nil
}

func GetMemoryMetrics() (MemoryMetrics, error) {
	stats, err := memory.Get()
	if err != nil {
		return MemoryMetrics{}, fmt.Errorf("collect memory metrics: %w", err)
	}

	return MemoryMetrics{
		TotalBytes:       stats.Total,
		UsedBytes:        stats.Used,
		FreeBytes:        stats.Free,
		AvailableBytes:   memoryAvailableBytes(stats, stats.Free),
		UsagePercent:     usagePercent(stats.Used, stats.Total),
		SwapTotalBytes:   stats.SwapTotal,
		SwapUsedBytes:    stats.SwapUsed,
		SwapFreeBytes:    stats.SwapFree,
		SwapUsagePercent: usagePercent(stats.SwapUsed, stats.SwapTotal),
	}, nil
}

func GetLoadMetrics() (LoadMetrics, error) {
	stats, err := loadavg.Get()
	if err != nil {
		return LoadMetrics{}, fmt.Errorf("collect load metrics: %w", err)
	}

	return LoadMetrics{
		Load1:  stats.Loadavg1,
		Load5:  stats.Loadavg5,
		Load15: stats.Loadavg15,
	}, nil
}

func GetGoRuntimeMetrics() GoRuntimeMetrics {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return GoRuntimeMetrics{
		GoroutineCount:     runtime.NumGoroutine(),
		GOMAXPROCS:         runtime.GOMAXPROCS(0),
		CgoCallCount:       runtime.NumCgoCall(),
		AllocBytes:         mem.Alloc,
		TotalAllocBytes:    mem.TotalAlloc,
		SysBytes:           mem.Sys,
		HeapAllocBytes:     mem.HeapAlloc,
		HeapSysBytes:       mem.HeapSys,
		HeapInuseBytes:     mem.HeapInuse,
		StackInuseBytes:    mem.StackInuse,
		NextGCBytes:        mem.NextGC,
		LastGCTimeUnixNano: mem.LastGC,
		NumGC:              mem.NumGC,
		PauseTotalNs:       mem.PauseTotalNs,
	}
}

func GetMachineMetrics() (MachineMetrics, error) {
	return GetMachineMetricsWithOptions(CollectOptions{})
}

func GetMachineMetricsWithOptions(options CollectOptions) (MachineMetrics, error) {
	diskPath := options.DiskPath
	if diskPath == "" {
		diskPath = DefaultDiskPath
	}

	cpuInterval := options.CPUSamplingInterval
	if cpuInterval == 0 {
		cpuInterval = DefaultCPUSamplingInterval
	}

	var metrics MachineMetrics
	var err error

	metrics.CPU, err = GetCPUMetricsWithInterval(cpuInterval)
	if err != nil {
		return metrics, err
	}
	metrics.Memory, err = GetMemoryMetrics()
	if err != nil {
		return metrics, err
	}
	metrics.Disk, err = GetDiskMetrics(diskPath)
	if err != nil {
		return metrics, err
	}
	metrics.Load, err = GetLoadMetrics()
	if err != nil {
		return metrics, err
	}
	metrics.GoRuntime = GetGoRuntimeMetrics()

	return metrics, nil
}

func usagePercent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(used) / float64(total) * 100
}

func memoryAvailableBytes(stats any, fallback uint64) uint64 {
	value := reflect.ValueOf(stats)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fallback
	}

	field := value.FieldByName("Available")
	if !field.IsValid() || field.Kind() != reflect.Uint64 {
		return fallback
	}
	return field.Uint()
}

func cpuUsagePercent(beforeTotal, beforeIdle, afterTotal, afterIdle uint64) float64 {
	totalDelta := afterTotal - beforeTotal
	if afterTotal < beforeTotal || totalDelta == 0 {
		return 0
	}

	var idleDelta uint64
	if afterIdle >= beforeIdle {
		idleDelta = afterIdle - beforeIdle
	}
	if idleDelta > totalDelta {
		return 0
	}

	return usagePercent(totalDelta-idleDelta, totalDelta)
}
