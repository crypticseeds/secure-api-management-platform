package models

import (
	"gorm.io/gorm"
)

type APIKey struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Key         string `gorm:"unique;not null"`
	LastUsedAt  *int64
	Description string
	User        User `gorm:"foreignKey:UserID"`
}
