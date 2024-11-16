package model

import (
	"time"
)

type Listing struct {
	ListingID    uint   `gorm:"primaryKey"`
	Title        string  `gorm:"size:2048;not null"`
	Price        float64 `gorm:"not null"`
	Location     string  `gorm:"size:512"`
	Description  string  `gorm:"type:text"`
	Link         string  `gorm:"size:1048;unique;not null"`
	Seller       string  `gorm:"size:100"`
	City         string  `gorm:"size:100"`
	Neighborhood string  `gorm:"size:100"`
	Meterage     int     `gorm:"not null"`
	Bedrooms     int     `gorm:"not null"`
	AdType       string  `gorm:"size:50"` // e.g., "فروش", "اجاره"
	Age          string  `gorm:"size:50"` // Age as string to capture different formats if needed
	HouseType    string  `gorm:"size:50"`
	Floor        int     `gorm:"not null"`
	Warehouse    bool    `gorm:"not null"`
	Elevator     bool    `gorm:"not null"`
	AdCreateDate string  `gorm:"size:50"` // Keeping as string as per your data example
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Images       []string   `gorm:"-"`    // Placeholder for associated images
}
