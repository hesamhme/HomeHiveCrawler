package model

import (
	"time"
)

type User struct {
	TelegramID   int        `gorm:"primaryKey"` // Use TelegramID as the primary key
	Username     string     `gorm:"size:100;not null;unique"`
	Role         string     `gorm:"size:20;not null"`                  // e.g., "admin", "user", "superadmin"
	Status       string     `gorm:"size:20;not null;default:'active'"` // active, inactive, suspended
	LastLogin    *time.Time `gorm:"null"`
	PasswordHash string     `gorm:"size:255;not null"` // Hash for password storage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
