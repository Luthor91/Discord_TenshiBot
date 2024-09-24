package models

import (
	"time"

	"gorm.io/gorm"
)

// UserShopCooldown représente le délai d'achat pour un utilisateur pour un article spécifique
type UserShopCooldown struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"` // Identifiant unique
	UserDiscordID string    `gorm:"not null"`   // ID Discord de l'utilisateur
	ItemID        uint      `gorm:"not null"`   // ID de l'article
	NextPurchase  time.Time // Prochain moment d'achat
}
