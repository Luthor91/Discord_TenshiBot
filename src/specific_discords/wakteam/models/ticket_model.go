package models

import (
	"time"
)

// Ticket représente une requête d'utilisateur
type Ticket struct {
	ID        uint      `gorm:"primaryKey"`
	Author    string    `gorm:"not null"` // Pseudo de l'auteur
	Content   string    `gorm:"not null"` // Contenu du message
	Status    string    `gorm:"not null"` // "en_cours", "valide", "refuse"
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
