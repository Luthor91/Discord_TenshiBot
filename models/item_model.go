package models

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"` // Clé primaire auto-incrémentée
	Name          string
	Quantity      int
	UserDiscordID string // Changez ceci en string pour correspondre à UserDiscordID
}
