package models

import (
	"time"

	"gorm.io/gorm"
)

// User représente un utilisateur avec des articles associés
type User struct {
	gorm.Model
	UserDiscordID   string `gorm:"uniqueIndex"` // Index unique
	Username        string
	TimeoutEnd      time.Time
}
