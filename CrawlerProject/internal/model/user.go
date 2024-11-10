package model

import (
	"time"

	"github.com/google/uuid" // Import for UUID support
)

type User struct {
	UserID       uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username     string     `gorm:"size:100;not null;unique"`
	TelegramID   *int       `gorm:"null;unique"`
	Role         string     `gorm:"size:20;not null"`                  // e.g., "admin", "user", "superadmin"
	Status       string     `gorm:"size:20;not null;default:'active'"` // active, inactive, suspended
	LastLogin    *time.Time `gorm:"null"`
	PasswordHash string     `gorm:"size:255;not null"` // Hash for password storage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
