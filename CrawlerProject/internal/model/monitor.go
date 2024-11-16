package model

import (
	"sync"
	"time"
)

type GoroutineMonitor struct {
	Stats    map[int64]*GoroutineStats
	StatsMux sync.RWMutex
	Done     chan struct{}
}

type GoroutineStats struct {
	GoroutineID    int64       `json:"goroutine_id"`
	StartTime      time.Time   `json:"start_time"`
	EndTime        time.Time   `json:"end_time"`
	URL            string      `json:"url"`
	City           string      `json:"city"`
	Type           string      `json:"type"`
	NumAdsFound    int         `json:"num_ads_found"`
	CPUUsage       []CPUSample `json:"cpu_usage"`
	MemoryUsage    []MemSample `json:"memory_usage"`
	PeakMemoryUsed uint64      `json:"peak_memory_used"`
	AvgCPUUsed     float64     `json:"avg_cpu_used"`
	Duration       float64     `json:"duration_seconds"`
}

type CPUSample struct {
	Timestamp time.Time `json:"timestamp"`
	Usage     float64   `json:"usage"`
}

type MemSample struct {
	Timestamp time.Time `json:"timestamp"`
	UsageKB   uint64    `json:"usage_kb"`
}
