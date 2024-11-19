package model

import (
	"time"
)

type Bookmark struct {
	BookmarkID uint      `gorm:"primaryKey"`
	UserID     int64     `gorm:"not null"` // TelegramID used to relate to User
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ListingID  uint      `gorm:"not null"` // Referencing the Listing model
	Listing    Listing   `gorm:"foreignKey:ListingID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
}
