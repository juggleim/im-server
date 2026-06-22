package commonservices

import (
	"im-server/commons/configures"
	"im-server/commons/metrics"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/logs"
	"time"
)

const (
	performanceMetricsCollectInterval = time.Minute
	performanceMetricsCleanupInterval = time.Hour
	performanceMetricsRetention       = 24 * time.Hour
)

var performanceMetricsTimer *time.Ticker
var performanceMetricsStartTimer *time.Timer
var performanceMetricsStop chan struct{}
var lastPerformanceMetricsCleanupTime time.Time

type PerformanceMetricsSnapshot struct {
	metrics.MachineMetrics
	ClientConnect ClientConnectMetrics `json:"client_connect"`
}

func CollectPerformanceMetricsSnapshot() (PerformanceMetricsSnapshot, error) {
	machineMetrics, err := metrics.GetMachineMetrics()
	if err != nil {
		return PerformanceMetricsSnapshot{}, err
	}
	return PerformanceMetricsSnapshot{
		MachineMetrics: machineMetrics,
		ClientConnect:  GetClientConnectMetrics(),
	}, nil
}

func FlattenPerformanceMetricRows(nodeName string, collectTime int64, snapshot PerformanceMetricsSnapshot) []dbs.PerformanceMetricDao {
	rows := make([]dbs.PerformanceMetricDao, 0, 33)
	add := func(metricType string, metricValue float64) {
		rows = append(rows, dbs.PerformanceMetricDao{
			NodeName:    nodeName,
			CollectTime: collectTime,
			MetricType:  metricType,
			MetricValue: metricValue,
		})
	}

	add("cpu.usage_percent", snapshot.CPU.UsagePercent)

	add("memory.total_bytes", float64(snapshot.Memory.TotalBytes))
	add("memory.used_bytes", float64(snapshot.Memory.UsedBytes))
	add("memory.free_bytes", float64(snapshot.Memory.FreeBytes))
	add("memory.available_bytes", float64(snapshot.Memory.AvailableBytes))
	add("memory.usage_percent", snapshot.Memory.UsagePercent)
	add("memory.swap_total_bytes", float64(snapshot.Memory.SwapTotalBytes))
	add("memory.swap_used_bytes", float64(snapshot.Memory.SwapUsedBytes))
	add("memory.swap_free_bytes", float64(snapshot.Memory.SwapFreeBytes))
	add("memory.swap_usage_percent", snapshot.Memory.SwapUsagePercent)

	add("disk.total_bytes", float64(snapshot.Disk.TotalBytes))
	add("disk.used_bytes", float64(snapshot.Disk.UsedBytes))
	add("disk.free_bytes", float64(snapshot.Disk.FreeBytes))
	add("disk.usage_percent", snapshot.Disk.UsagePercent)

	add("load.load1", snapshot.Load.Load1)
	add("load.load5", snapshot.Load.Load5)
	add("load.load15", snapshot.Load.Load15)

	add("go_runtime.goroutine_count", float64(snapshot.GoRuntime.GoroutineCount))
	add("go_runtime.gomaxprocs", float64(snapshot.GoRuntime.GOMAXPROCS))
	add("go_runtime.cgo_call_count", float64(snapshot.GoRuntime.CgoCallCount))
	add("go_runtime.alloc_bytes", float64(snapshot.GoRuntime.AllocBytes))
	add("go_runtime.total_alloc_bytes", float64(snapshot.GoRuntime.TotalAllocBytes))
	add("go_runtime.sys_bytes", float64(snapshot.GoRuntime.SysBytes))
	add("go_runtime.heap_alloc_bytes", float64(snapshot.GoRuntime.HeapAllocBytes))
	add("go_runtime.heap_sys_bytes", float64(snapshot.GoRuntime.HeapSysBytes))
	add("go_runtime.heap_inuse_bytes", float64(snapshot.GoRuntime.HeapInuseBytes))
	add("go_runtime.stack_inuse_bytes", float64(snapshot.GoRuntime.StackInuseBytes))
	add("go_runtime.next_gc_bytes", float64(snapshot.GoRuntime.NextGCBytes))
	add("go_runtime.last_gc_time_unix_nano", float64(snapshot.GoRuntime.LastGCTimeUnixNano))
	add("go_runtime.num_gc", float64(snapshot.GoRuntime.NumGC))
	add("go_runtime.pause_total_ns", float64(snapshot.GoRuntime.PauseTotalNs))

	add("client_connect.online_user_count", float64(snapshot.ClientConnect.OnlineUserCount))
	add("client_connect.user_connect_count", float64(snapshot.ClientConnect.UserConnectCount))
	add("client_connect.session_connect_count", float64(snapshot.ClientConnect.SessionConnectCount))

	return rows
}

func StartPerformanceMetricsCollect() {
	StopPerformanceMetricsCollect()
	if !configures.Config.PerformanceMetrics.IsOpen {
		return
	}

	performanceMetricsStop = make(chan struct{})
	nextCollectTime := nextAlignedCollectTime(time.Now(), performanceMetricsCollectInterval)
	performanceMetricsStartTimer = time.NewTimer(time.Until(nextCollectTime))
	go func(startTimer *time.Timer, stop <-chan struct{}) {
		select {
		case <-startTimer.C:
			collectAndPersistPerformanceMetrics(nextCollectTime)
		case <-stop:
			return
		}

		performanceMetricsTimer = time.NewTicker(performanceMetricsCollectInterval)
		defer performanceMetricsTimer.Stop()
		for {
			select {
			case task := <-performanceMetricsTimer.C:
				collectAndPersistPerformanceMetrics(alignedCollectTime(task, performanceMetricsCollectInterval))
			case <-stop:
				return
			}
		}
	}(performanceMetricsStartTimer, performanceMetricsStop)
}

func StopPerformanceMetricsCollect() {
	if performanceMetricsStartTimer != nil {
		performanceMetricsStartTimer.Stop()
		performanceMetricsStartTimer = nil
	}
	if performanceMetricsTimer != nil {
		performanceMetricsTimer.Stop()
		performanceMetricsTimer = nil
	}
	if performanceMetricsStop != nil {
		close(performanceMetricsStop)
		performanceMetricsStop = nil
	}
}

func collectAndPersistPerformanceMetrics(collectTime time.Time) {
	snapshot, err := CollectPerformanceMetricsSnapshot()
	if err != nil {
		logs.NewLogEntity().Errorf("collect performance metrics failed:%v", err)
		return
	}

	rows := FlattenPerformanceMetricRows(configures.Config.NodeName, collectTime.UnixMilli(), snapshot)
	dao := dbs.PerformanceMetricDao{}
	if err := dao.BatchInsert(rows); err != nil {
		logs.NewLogEntity().Errorf("persist performance metrics failed:%v", err)
		return
	}
	if !shouldCleanupPerformanceMetrics(collectTime) {
		return
	}
	if err := dao.DeleteBeforeCollectTime(configures.Config.NodeName, performanceMetricsRetentionCutoff(collectTime)); err != nil {
		logs.NewLogEntity().Errorf("cleanup performance metrics failed:%v", err)
	}
}

func shouldCleanupPerformanceMetrics(collectTime time.Time) bool {
	if lastPerformanceMetricsCleanupTime.IsZero() || collectTime.Sub(lastPerformanceMetricsCleanupTime) >= performanceMetricsCleanupInterval {
		lastPerformanceMetricsCleanupTime = collectTime
		return true
	}
	return false
}

func performanceMetricsRetentionCutoff(collectTime time.Time) int64 {
	return collectTime.Add(-performanceMetricsRetention).UnixMilli()
}

func alignedCollectTime(t time.Time, interval time.Duration) time.Time {
	if interval <= 0 {
		return t
	}
	intervalMs := int64(interval / time.Millisecond)
	if intervalMs <= 0 {
		return t
	}
	timeMs := t.UnixMilli()
	return time.UnixMilli(timeMs / intervalMs * intervalMs)
}

func nextAlignedCollectTime(t time.Time, interval time.Duration) time.Time {
	aligned := alignedCollectTime(t, interval)
	if !aligned.After(t) {
		aligned = aligned.Add(interval)
	}
	return aligned
}
