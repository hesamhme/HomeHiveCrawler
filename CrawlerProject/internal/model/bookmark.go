package model

import (
	"github.com/google/uuid"
	"time"
)

type Bookmark struct {
	BookmarkID uint      `gorm:"primaryKey"`
	UserID     uuid.UUID `gorm:"not null"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ListingID  uint      `gorm:"not null"`
	Listing    Listing   `gorm:"foreignKey:ListingID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	Notes      string `gorm:"type:text;null"` // Optional notes for user
}
