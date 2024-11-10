package model

import (
	"time"
)

type Listing struct {
	ListingID    uint   `gorm:"primaryKey"`
	Title        string `gorm:"size:255;not null"`
	Description  string `gorm:"type:text"`
	Price        float64
	City         string `gorm:"size:100"`
	Neighborhood string `gorm:"size:100"`
	Area         float64
	Rooms        int
	Status       string `gorm:"size:20"` // e.g., "rent", "sell"
	Floor        *string
	HasStorage   bool
	HasElevator  bool
	Source       string     `gorm:"size:100"` // e.g., "divar", "sheypoor"
	URL          string     `gorm:"type:text"`
	ExpiresAt    *time.Time `gorm:"null"` // Expiration date of the listing
	Images       []string   `gorm:"-"`    // Placeholder for associated images
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
