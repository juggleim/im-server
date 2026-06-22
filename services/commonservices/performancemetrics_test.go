package commonservices

import (
	"im-server/commons/configures"
	"im-server/commons/metrics"
	"testing"
	"time"
)

func TestFlattenPerformanceMetricRows(t *testing.T) {
	snapshot := PerformanceMetricsSnapshot{
		MachineMetrics: metrics.MachineMetrics{
			CPU: metrics.CPUMetrics{
				UsagePercent: 12.5,
			},
			Memory: metrics.MemoryMetrics{
				TotalBytes:       100,
				UsedBytes:        40,
				FreeBytes:        60,
				AvailableBytes:   55,
				UsagePercent:     40,
				SwapTotalBytes:   10,
				SwapUsedBytes:    2,
				SwapFreeBytes:    8,
				SwapUsagePercent: 20,
			},
			Disk: metrics.DiskMetrics{
				TotalBytes:   200,
				UsedBytes:    50,
				FreeBytes:    150,
				UsagePercent: 25,
			},
			Load: metrics.LoadMetrics{
				Load1:  1,
				Load5:  2,
				Load15: 3,
			},
			GoRuntime: metrics.GoRuntimeMetrics{
				GoroutineCount:     4,
				GOMAXPROCS:         5,
				CgoCallCount:       6,
				AllocBytes:         7,
				TotalAllocBytes:    8,
				SysBytes:           9,
				HeapAllocBytes:     10,
				HeapSysBytes:       11,
				HeapInuseBytes:     12,
				StackInuseBytes:    13,
				NextGCBytes:        14,
				LastGCTimeUnixNano: 15,
				NumGC:              16,
				PauseTotalNs:       17,
			},
		},
		ClientConnect: ClientConnectMetrics{
			OnlineUserCount:     18,
			UserConnectCount:    19,
			SessionConnectCount: 20,
		},
	}

	rows := FlattenPerformanceMetricRows("node-a", 123456, snapshot)
	metricMap := map[string]float64{}
	for _, row := range rows {
		if row.NodeName != "node-a" {
			t.Fatalf("NodeName = %q, want node-a", row.NodeName)
		}
		if row.CollectTime != 123456 {
			t.Fatalf("CollectTime = %d, want 123456", row.CollectTime)
		}
		metricMap[row.MetricType] = row.MetricValue
	}

	expected := map[string]float64{
		"cpu.usage_percent":                    12.5,
		"memory.usage_percent":                 40,
		"disk.used_bytes":                      50,
		"load.load1":                           1,
		"go_runtime.goroutine_count":           4,
		"go_runtime.heap_sys_bytes":            11,
		"client_connect.online_user_count":     18,
		"client_connect.user_connect_count":    19,
		"client_connect.session_connect_count": 20,
	}
	for metricType, want := range expected {
		if got, ok := metricMap[metricType]; !ok {
			t.Fatalf("missing metric type %s", metricType)
		} else if got != want {
			t.Fatalf("metric %s = %v, want %v", metricType, got, want)
		}
	}
}

func TestStartPerformanceMetricsCollectDisabled(t *testing.T) {
	configures.Config.PerformanceMetrics.IsOpen = false
	StartPerformanceMetricsCollect()
	defer StopPerformanceMetricsCollect()

	if performanceMetricsTimer != nil {
		t.Fatal("performanceMetricsTimer is not nil when persistence is disabled")
	}
}

func TestAlignedCollectTime(t *testing.T) {
	interval := time.Minute
	input := time.Date(2026, 6, 17, 10, 7, 42, 123*int(time.Millisecond), time.UTC)
	want := time.Date(2026, 6, 17, 10, 7, 0, 0, time.UTC)

	if got := alignedCollectTime(input, interval); !got.Equal(want) {
		t.Fatalf("alignedCollectTime() = %s, want %s", got, want)
	}
}

func TestPerformanceMetricsRetentionCutoff(t *testing.T) {
	collectTime := time.Date(2026, 6, 17, 10, 5, 0, 0, time.UTC)
	want := time.Date(2026, 6, 16, 10, 5, 0, 0, time.UTC).UnixMilli()

	if got := performanceMetricsRetentionCutoff(collectTime); got != want {
		t.Fatalf("performanceMetricsRetentionCutoff() = %d, want %d", got, want)
	}
}

func TestShouldCleanupPerformanceMetrics(t *testing.T) {
	lastPerformanceMetricsCleanupTime = time.Time{}
	defer func() {
		lastPerformanceMetricsCleanupTime = time.Time{}
	}()

	first := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	if !shouldCleanupPerformanceMetrics(first) {
		t.Fatal("first cleanup should be allowed")
	}
	if shouldCleanupPerformanceMetrics(first.Add(55 * time.Minute)) {
		t.Fatal("cleanup inside one hour should be skipped")
	}
	if !shouldCleanupPerformanceMetrics(first.Add(time.Hour)) {
		t.Fatal("cleanup after one hour should be allowed")
	}
}

func TestNextAlignedCollectTime(t *testing.T) {
	interval := time.Minute
	tests := []struct {
		name string
		now  time.Time
		want time.Time
	}{
		{
			name: "between buckets",
			now:  time.Date(2026, 6, 17, 10, 7, 42, 0, time.UTC),
			want: time.Date(2026, 6, 17, 10, 8, 0, 0, time.UTC),
		},
		{
			name: "on bucket boundary",
			now:  time.Date(2026, 6, 17, 10, 10, 0, 0, time.UTC),
			want: time.Date(2026, 6, 17, 10, 11, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nextAlignedCollectTime(tt.now, interval); !got.Equal(tt.want) {
				t.Fatalf("nextAlignedCollectTime() = %s, want %s", got, tt.want)
			}
		})
	}
}
