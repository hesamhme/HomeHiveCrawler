package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
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

// StartTracking begins monitoring a new goroutine
func (gm *GoroutineMonitor) StartTracking(city, _type string) *GoroutineStats {
	stats := &GoroutineStats{
		GoroutineID: time.Now().UnixNano(), // Use timestamp as unique ID
		StartTime:   time.Now(),
		City:        city,
		Type:        _type,
		CPUUsage:    make([]CPUSample, 0),
		MemoryUsage: make([]MemSample, 0),
	}

	gm.StatsMux.Lock()
	gm.Stats[stats.GoroutineID] = stats
	gm.StatsMux.Unlock()

	// Start monitoring goroutine resources
	go gm.monitorResources(stats.GoroutineID)

	return stats
}

// StopTracking ends monitoring for a goroutine
func (gm *GoroutineMonitor) StopTracking(goroutineID int64) {
	gm.StatsMux.Lock()
	if stats, exists := gm.Stats[goroutineID]; exists {
		stats.EndTime = time.Now()
		stats.Duration = stats.EndTime.Sub(stats.StartTime).Seconds()

		// Calculate average CPU usage
		var totalCPU float64
		for _, sample := range stats.CPUUsage {
			totalCPU += sample.Usage
		}
		if len(stats.CPUUsage) > 0 {
			stats.AvgCPUUsed = totalCPU / float64(len(stats.CPUUsage))
		}
	}
	gm.StatsMux.Unlock()
}

// monitorResources continuously monitors resource usage for a goroutine
func (gm *GoroutineMonitor) monitorResources(goroutineID int64) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			gm.StatsMux.Lock()
			if stats, exists := gm.Stats[goroutineID]; exists {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				// Record memory sample
				memSample := MemSample{
					Timestamp: time.Now(),
					UsageKB:   m.Alloc / 1024,
				}
				stats.MemoryUsage = append(stats.MemoryUsage, memSample)

				// Update peak memory if current is higher
				if m.Alloc > stats.PeakMemoryUsed {
					stats.PeakMemoryUsed = m.Alloc
				}

				// Record CPU sample (simplified version)
				cpuSample := CPUSample{
					Timestamp: time.Now(),
					Usage:     getCPUUsage(), // You'll need to implement this
				}
				stats.CPUUsage = append(stats.CPUUsage, cpuSample)
			}
			gm.StatsMux.Unlock()

		case <-gm.Done:
			return
		}
	}
}

// SaveStats saves the monitoring statistics to a file
func (gm *GoroutineMonitor) SaveStats(outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(outputDir,
		fmt.Sprintf("goroutine_stats_%s.json", time.Now().Format("2006-01-02_15-04-05")))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create stats file: %w", err)
	}
	defer file.Close()

	gm.StatsMux.RLock()
	defer gm.StatsMux.RUnlock()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(gm.Stats); err != nil {
		return fmt.Errorf("failed to encode stats: %w", err)
	}

	log.Printf("Goroutine statistics saved to %s", filename)
	return nil
}

// Helper function to get CPU usage (simplified version)
func getCPUUsage() float64 {
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return 0
	}

	cpu, err := p.CPUPercent()
	if err != nil {
		return 0
	}

	return cpu
}
func NewGoroutineMonitor() *GoroutineMonitor {
	return &GoroutineMonitor{
		Stats: make(map[int64]*GoroutineStats),
		Done:  make(chan struct{}),
	}

}
