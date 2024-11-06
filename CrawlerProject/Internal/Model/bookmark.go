package Model

import (
	"time"
)

type Bookmark struct {
	BookmarkID uint    `gorm:"primaryKey"`
	UserID     uint    `gorm:"not null"`
	User       User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ListingID  uint    `gorm:"not null"`
	Listing    Listing `gorm:"foreignKey:ListingID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
}
