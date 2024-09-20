package models

import (
	"time"

	"gorm.io/gorm"
)

// Warn représente un avertissement lié à un utilisateur.
type Warn struct {
	gorm.Model
	UserDiscordID string    // Lien avec l'utilisateur via son DiscordID
	Reason        string    // Raison du warn
	WarnedAt      time.Time // Date du warn
	AdminID       string    // ID de l'admin/modérateur ayant émis le warn
}
