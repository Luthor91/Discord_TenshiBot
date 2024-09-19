package models

import (
	"time"

	"gorm.io/gorm"
)

// UserShopCooldown représente le délai d'achat pour un utilisateur pour un article spécifique
type UserShopCooldown struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"` // Identifiant unique
	UserDiscordID string    `gorm:"not null"`   // Clé étrangère vers User (ID Discord)
	ItemID        uint      `gorm:"not null"`   // Clé étrangère vers ShopItem
	NextPurchase  time.Time // Prochain moment où l'utilisateur pourra acheter l'article
	ShopItem      ShopItem  `gorm:"foreignKey:ItemID"` // Relation avec ShopItem
}
