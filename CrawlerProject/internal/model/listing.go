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
	Area         float64    `gorm:"column:meterage" json:"meterage"`
	Rooms        int        `gorm:"column:bedrooms" json:"bedrooms"`
	Status       string     `gorm:"size:20" json:"adType"`
	Floor        int        `gorm:"" json:"floor"`
	HasStorage   bool       `gorm:"column:warehouse" json:"warehouse"`
	HasElevator  bool       `gorm:"column:elevator" json:"elevator"`
	HasParking   bool       `gorm:"column:parking" json:"parking"`
	Source       string     `gorm:"size:100" json:"source"`
	URL          string     `gorm:"type:text" json:"link"`
	Seller       string     `gorm:"size:100" json:"seller"`
	HouseType    string     `gorm:"size:50" json:"houseType"`
	Age          string     `gorm:"size:50" json:"age"`
	ExpiresAt    *time.Time `gorm:"null" json:"expiresAt"`
	Images       []string   `gorm:"-" json:"images"`
	CreatedAt    time.Time  `gorm:"" json:"adCreateDate"`
	UpdatedAt    time.Time  `gorm:"" json:"updatedAt"`
}
