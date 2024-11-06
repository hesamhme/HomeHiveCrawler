package Model

import (
	"time"
)

type AdminLog struct {
	LogID       uint   `gorm:"primaryKey"`
	AdminID     uint   `gorm:"not null"`
	Admin       User   `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE"`
	Action      string `gorm:"size:100"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
}
