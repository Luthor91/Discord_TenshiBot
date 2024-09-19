package models

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"` // Clé primaire auto-incrémentée
	Name          string
	Quantity      int
	UserDiscordID uint // Ajout de la clé étrangère pour la relation avec User
}
