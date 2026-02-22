package pkg_logic

import (
	"os"
	"runtime"
	"time"
)

func GetClockState() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"local_time":  now.Format(time.RFC3339),
		"local_date":  now.Format("2006-01-02"),
		"day_of_week": now.Weekday().String(),
		"timestamp":   now.Unix(),
	}
}

func GetCalendarState() map[string]interface{} {
	now := time.Now()
	return map[string]interface{}{
		"month": int(now.Month()),
		"day":   now.Day(),
		"year":  now.Year(),
	}
}

func GetMetricsState() map[string]interface{} {
	hostname, _ := os.Hostname()
	return map[string]interface{}{
		"cpu_cores":  runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
		"os":         runtime.GOOS,
		"hostname":   hostname,
	}
}