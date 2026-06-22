//go:build linux || darwin || freebsd || netbsd || openbsd || windows

package metrics

import "testing"

func TestGetDiskMetrics(t *testing.T) {
	metrics, err := GetDiskMetrics(t.TempDir())
	if err != nil {
		t.Fatalf("GetDiskMetrics() error = %v", err)
	}

	if metrics.Path == "" {
		t.Fatal("Path is empty")
	}
	if metrics.TotalBytes == 0 {
		t.Fatal("TotalBytes = 0, want > 0")
	}
	if metrics.UsedBytes > metrics.TotalBytes {
		t.Fatalf("UsedBytes = %d, want <= TotalBytes %d", metrics.UsedBytes, metrics.TotalBytes)
	}
	if metrics.FreeBytes > metrics.TotalBytes {
		t.Fatalf("FreeBytes = %d, want <= TotalBytes %d", metrics.FreeBytes, metrics.TotalBytes)
	}
	if metrics.UsagePercent < 0 || metrics.UsagePercent > 100 {
		t.Fatalf("UsagePercent = %v, want between 0 and 100", metrics.UsagePercent)
	}
}
