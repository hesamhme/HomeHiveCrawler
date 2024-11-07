package Model

import (
	"time"
)

type CrawlerLog struct {
	LogID        uint   `gorm:"primaryKey"`
	CrawlerName  string `gorm:"size:100"`
	StartTime    time.Time
	EndTime      time.Time
	CPUUsage     float64
	MemoryUsage  float64
	Status       string `gorm:"size:20"`
	ErrorMessage string `gorm:"type:text"`
	CreatedAt    time.Time
}