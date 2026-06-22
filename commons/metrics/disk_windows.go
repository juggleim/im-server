//go:build windows

package metrics

import (
	"fmt"
	"syscall"
	"unsafe"
)

func GetDiskMetrics(path string) (DiskMetrics, error) {
	if path == "" {
		return DiskMetrics{}, ErrInvalidDiskPath
	}

	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return DiskMetrics{}, fmt.Errorf("collect disk metrics for %q: %w", path, err)
	}

	var freeBytesAvailableToCaller uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GetDiskFreeSpaceExW")
	ret, _, callErr := proc.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailableToCaller)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if ret == 0 {
		return DiskMetrics{}, fmt.Errorf("collect disk metrics for %q: %w", path, callErr)
	}

	used := totalNumberOfBytes - totalNumberOfFreeBytes
	return DiskMetrics{
		Path:         path,
		TotalBytes:   totalNumberOfBytes,
		UsedBytes:    used,
		FreeBytes:    totalNumberOfFreeBytes,
		UsagePercent: usagePercent(used, totalNumberOfBytes),
	}, nil
}
