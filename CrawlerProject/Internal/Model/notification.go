package Model

import (
	"time"
)

type Notification struct {
	NotificationID uint    `gorm:"primaryKey"`
	UserID         uint    `gorm:"not null"`
	User           User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FilterID       uint    `gorm:"not null"`
	Filter         Filter  `gorm:"foreignKey:FilterID;constraint:OnDelete:CASCADE"`
	ListingID      uint    `gorm:"not null"`
	Listing        Listing `gorm:"foreignKey:ListingID;constraint:OnDelete:CASCADE"`
	Status         string  `gorm:"size:20"`
	CreatedAt      time.Time
	SentAt         *time.Time
}
