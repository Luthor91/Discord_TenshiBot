package models

import "time"

// Investment représente un investissement effectué par un utilisateur.
type Investment struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      // Référence à l'utilisateur
	Amount     int       // Montant investi
	CreatedAt  time.Time // Date de l'investissement
	InvestedAt time.Time // Date à laquelle l'utilisateur a récupéré l'argent
	Result     int       // Gain ou perte après traitement
	Status     string    // "pending", "success", "fail"
	Factor     float64   // Facteur multiplicateur pour le gain ou la perte
	Cooldown   time.Time // Timestamp du cooldown de 24 heures
}
