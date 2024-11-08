package Model

import (
	"github.com/google/uuid"
	"time"
)

type AdminLog struct {
	LogID       uint      `gorm:"primaryKey"`
	AdminID     uuid.UUID `gorm:"not null"`
	Admin       User      `gorm:"foreignKey:AdminID;constraint:OnDelete:CASCADE"`
	Action      string    `gorm:"size:100"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
}
