package tools

import (
	"time"

	"golang.org/x/time/rate"
)

// NewPerSecondLimiter 创建一个以秒为单位的限频器，每秒最多允许 limit 次通过
// limit: 每秒允许的最大请求数
func NewPerSecondLimiter(limit int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(limit), limit)
}

// NewPerMinuteLimiter 创建一个以分钟为单位的限频器，每分钟最多允许 limit 次通过
// limit: 每分钟允许的最大请求数
func NewPerMinuteLimiter(limit int) *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Minute/time.Duration(limit)), limit)
}

// NewPerHourLimiter 创建一个以小时为单位的限频器，每小时最多允许 limit 次通过
// limit: 每小时允许的最大请求数
func NewPerHourLimiter(limit int) *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Hour/time.Duration(limit)), limit)
}
