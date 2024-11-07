package Model

import (
	"time"
)

type SearchHistory struct {
	SearchID    uint   `gorm:"primaryKey"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	SearchTerms string `gorm:"type:text"`
	CreatedAt   time.Time
}
