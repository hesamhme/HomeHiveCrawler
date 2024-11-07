package Model

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
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Source       string `gorm:"size:100"` // e.g., "divar", "sheypoor"
	URL          string `gorm:"type:text"`
}
