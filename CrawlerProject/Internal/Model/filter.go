package Model

import (
	"github.com/google/uuid"
	"time"
)

type Filter struct {
	FilterID     uint      `gorm:"primaryKey"`
	UserID       uuid.UUID `gorm:"not null"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	PriceMin     float64
	PriceMax     float64
	City         string `gorm:"size:100"`
	Neighborhood string `gorm:"size:100"`
	AreaMin      float64
	AreaMax      float64
	RoomsMin     int
	RoomsMax     int
	Status       string `gorm:"size:20"` // e.g., "rent", "buy"
	FloorMin     *int
	FloorMax     *int
	HasStorage   bool
	HasElevator  bool
	Radius       float64 `gorm:"-"` // Radius for map-based filtering
	CreatedAt    time.Time
}
