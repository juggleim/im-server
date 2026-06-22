//go:build !linux && !darwin && !freebsd && !netbsd && !openbsd && !windows

package metrics

import (
	"fmt"
	"runtime"
)

func GetDiskMetrics(path string) (DiskMetrics, error) {
	if path == "" {
		return DiskMetrics{}, ErrInvalidDiskPath
	}
	return DiskMetrics{}, fmt.Errorf("collect disk metrics for %q: unsupported platform %s", path, runtime.GOOS)
}
