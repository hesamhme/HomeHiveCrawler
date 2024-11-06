package Model

import (
	"time"
)

type User struct {
	UserID     uint   `gorm:"primaryKey"`
	Username   string `gorm:"size:100;not null;unique"`
	TelegramID *int   `gorm:"null"`
	Role       string `gorm:"size:20;not null"` // e.g., "admin", "user", "superadmin"
	Email      string `gorm:"size:100;unique"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
