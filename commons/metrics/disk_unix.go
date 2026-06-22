//go:build linux || darwin || freebsd || netbsd || openbsd

package metrics

import (
	"fmt"
	"syscall"
)

func GetDiskMetrics(path string) (DiskMetrics, error) {
	if path == "" {
		return DiskMetrics{}, ErrInvalidDiskPath
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return DiskMetrics{}, fmt.Errorf("collect disk metrics for %q: %w", path, err)
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	return DiskMetrics{
		Path:         path,
		TotalBytes:   total,
		UsedBytes:    used,
		FreeBytes:    free,
		UsagePercent: usagePercent(used, total),
	}, nil
}
