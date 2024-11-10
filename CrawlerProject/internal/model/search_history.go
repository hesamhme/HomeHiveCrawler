package model

import (
	"github.com/google/uuid"
	"time"
)

type SearchHistory struct {
	SearchID    uint      `gorm:"primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SearchTerms string    `gorm:"type:text"`
	CreatedAt   time.Time
}
