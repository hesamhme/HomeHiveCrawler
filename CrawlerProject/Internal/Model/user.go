package Model

import (
	"time"
)

type User struct {
	UserID     uint   `gorm:"primaryKey"`
	Username   string `gorm:"size:100;not null;unique"`
	TelegramID *int   `gorm:"null;unique"`
	Role       string `gorm:"size:20;not null"` // e.g., "admin", "user", "superadmin"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
