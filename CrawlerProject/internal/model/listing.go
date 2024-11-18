package model

import (
	"time"
)

type Listing struct {
	ListingID    uint       `gorm:"primaryKey" json:"listingId"`
	Title        string     `gorm:"size:255;not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Price        float64    `gorm:"not null" json:"price"`
	City         string     `gorm:"size:100" json:"city"`
	Neighborhood string     `gorm:"size:100" json:"neighborhood"`
	Area         float64    `gorm:"column:meterage" json:"meterage"`   // Mapped from Meterage
	Rooms        int        `gorm:"column:bedrooms" json:"bedrooms"`   // Mapped from Bedrooms
	Status       string     `gorm:"size:20" json:"adType"`             // Mapped from AdType
	Floor        int        `gorm:"" json:"floor"`                     // Changed from *string to int
	HasStorage   bool       `gorm:"column:warehouse" json:"warehouse"` // Mapped from WareHouse
	HasElevator  bool       `gorm:"column:elevator" json:"elevator"`   // Mapped from Elevator
	HasParking   bool       `gorm:"column:parking" json:"parking"`     // Mapped from Parking
	Source       string     `gorm:"size:100" json:"source"`            // e.g., "divar", "sheypoor"
	URL          string     `gorm:"type:text" json:"url"`              // Mapped from URL
	Seller       string     `gorm:"size:100" json:"seller"`            // Added from HouseAd
	HouseType    string     `gorm:"size:50" json:"houseType"`          // Added from HouseAd
	Age          string     `gorm:"size:50" json:"age"`                // Added from HouseAd
	ExpiresAt    *time.Time `gorm:"null" json:"expiresAt"`             // Kept from original Listing
	Images       []string   `gorm:"-" json:"images"`                   // Using json tag from HouseAd
	CreatedAt    time.Time  `gorm:"" json:"adCreateDate"`              // Mapped from AdCreateDate
	UpdatedAt    time.Time  `gorm:"" json:"updatedAt"`
}
