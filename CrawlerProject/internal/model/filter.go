package model

import (
	"time"
)

type Filter struct {
	FilterID        uint  `gorm:"primaryKey"`
	UserID          int64 `gorm:"not null"`
	User            User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PriceMin        float64
	PriceMax        float64
	City            string `gorm:"size:100"`
	Neighborhood    string `gorm:"size:100"`
	AreaMin         float64
	AreaMax         float64
	RoomsMin        int
	RoomsMax        int
	Status          string `gorm:"size:20"`
	BuildingAgeMin  int
	BuildingAgeMax  int
	PropertyType    string `gorm:"size:50"`
	FloorMin        *int
	FloorMax        *int
	HasStorage      bool
	HasElevator     bool
	CreationDateMin time.Time
	CreationDateMax time.Time
	Latitude        float64
	Longitude       float64
	Radius          float64
	CreatedAt       time.Time
}
