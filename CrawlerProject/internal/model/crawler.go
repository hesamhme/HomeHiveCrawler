package model

import (
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type Crawler struct {
	// Configuration
	Config CrawlerConfig

	// State management
	LastRunTime time.Time
	RunMutex    sync.Mutex

	// Resource monitoring
	// monitor *ResourceMonitor
	GoroutineMonitor *GoroutineMonitor

	// Concurrency control
	UrlSemaphore chan struct{}
	AdsSemaphore chan struct{}

	// Error handling
	ErrorChan chan error

	// Results management
	ResultsChan chan Listing
	AdsMutex    sync.Mutex
}

type CrawlerConfig struct {
	// Time configuration
	RunInterval        time.Duration
	MinTimeBetweenRuns time.Duration
	PageTimeout        time.Duration
	AdTimeout          time.Duration

	// Concurrency limits
	MaxURLConcurrency int
	MaxAdConcurrency  int

	// Target configuration
	Cities []string
	Types  []string

	// Output configuration
	OutputDir string

	// Browser configuration
	ChromeFlags []chromedp.ExecAllocatorOption
}

// type HouseAd struct {
// 	Title        string    `json:"title"`
// 	Price        uint64    `json:"price"`
// 	Location     string    `json:"location"`
// 	Description  string    `json:"description"`
// 	Link         string    `json:"link"`
// 	Seller       string    `json:"seller"`
// 	City         string    `json:"city"`
// 	Neighborhood string    `json:"neighborhood"`
// 	Meterage     int       `json:"meterage"`
// 	Bedrooms     int       `json:"bedrooms"`
// 	AdType       string    `json:"adType"`
// 	Age          string    `json:"age"`
// 	HouseType    string    `json:"houseType"`
// 	Floor        int       `json:"floor"`
// 	WareHouse    bool      `json:"warehouse"`
// 	Elevator     bool      `json:"elevator"`
// 	AdCreateDate time.Time `json:"adCreateDate"`
// 	Images       []string  `json:"images"`
// }
