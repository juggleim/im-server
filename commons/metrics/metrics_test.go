package metrics

import (
	"errors"
	"testing"
	"time"
)

func TestGetGoRuntimeMetrics(t *testing.T) {
	metrics := GetGoRuntimeMetrics()

	if metrics.GoroutineCount <= 0 {
		t.Fatalf("GoroutineCount = %d, want > 0", metrics.GoroutineCount)
	}
	if metrics.GOMAXPROCS <= 0 {
		t.Fatalf("GOMAXPROCS = %d, want > 0", metrics.GOMAXPROCS)
	}
	if metrics.AllocBytes == 0 {
		t.Fatal("AllocBytes = 0, want runtime memory stats")
	}
	if metrics.SysBytes == 0 {
		t.Fatal("SysBytes = 0, want runtime memory stats")
	}
	if metrics.HeapSysBytes == 0 {
		t.Fatal("HeapSysBytes = 0, want runtime heap stats")
	}
}

func TestUsagePercent(t *testing.T) {
	tests := []struct {
		name  string
		used  uint64
		total uint64
		want  float64
	}{
		{name: "zero total", used: 10, total: 0, want: 0},
		{name: "zero used", used: 0, total: 100, want: 0},
		{name: "half used", used: 50, total: 100, want: 50},
		{name: "all used", used: 100, total: 100, want: 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usagePercent(tt.used, tt.total); got != tt.want {
				t.Fatalf("usagePercent(%d, %d) = %v, want %v", tt.used, tt.total, got, tt.want)
			}
		})
	}
}

func TestMemoryAvailableBytes(t *testing.T) {
	t.Run("uses available field when present", func(t *testing.T) {
		stats := struct {
			Available uint64
		}{Available: 123}

		if got := memoryAvailableBytes(&stats, 456); got != 123 {
			t.Fatalf("memoryAvailableBytes() = %d, want 123", got)
		}
	})

	t.Run("falls back when available field is absent", func(t *testing.T) {
		stats := struct {
			Free uint64
		}{Free: 123}

		if got := memoryAvailableBytes(&stats, 456); got != 456 {
			t.Fatalf("memoryAvailableBytes() = %d, want 456", got)
		}
	})
}

func TestCPUUsagePercent(t *testing.T) {
	tests := []struct {
		name        string
		beforeTotal uint64
		beforeIdle  uint64
		afterTotal  uint64
		afterIdle   uint64
		want        float64
	}{
		{name: "zero delta", beforeTotal: 100, beforeIdle: 50, afterTotal: 100, afterIdle: 50, want: 0},
		{name: "half busy", beforeTotal: 100, beforeIdle: 50, afterTotal: 200, afterIdle: 100, want: 50},
		{name: "all busy", beforeTotal: 100, beforeIdle: 50, afterTotal: 200, afterIdle: 50, want: 100},
		{name: "counter rollback", beforeTotal: 200, beforeIdle: 50, afterTotal: 100, afterIdle: 70, want: 0},
		{name: "idle exceeds total delta", beforeTotal: 100, beforeIdle: 10, afterTotal: 150, afterIdle: 80, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cpuUsagePercent(tt.beforeTotal, tt.beforeIdle, tt.afterTotal, tt.afterIdle)
			if got != tt.want {
				t.Fatalf("cpuUsagePercent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCPUMetricsWithIntervalRejectsInvalidInterval(t *testing.T) {
	_, err := GetCPUMetricsWithInterval(0)
	if !errors.Is(err, ErrInvalidCPUSamplingInterval) {
		t.Fatalf("GetCPUMetricsWithInterval(0) error = %v, want %v", err, ErrInvalidCPUSamplingInterval)
	}

	_, err = GetCPUMetricsWithInterval(-time.Millisecond)
	if !errors.Is(err, ErrInvalidCPUSamplingInterval) {
		t.Fatalf("GetCPUMetricsWithInterval(-1ms) error = %v, want %v", err, ErrInvalidCPUSamplingInterval)
	}
}

func TestGetDiskMetricsRejectsEmptyPath(t *testing.T) {
	_, err := GetDiskMetrics("")
	if !errors.Is(err, ErrInvalidDiskPath) {
		t.Fatalf("GetDiskMetrics(empty) error = %v, want %v", err, ErrInvalidDiskPath)
	}
}
