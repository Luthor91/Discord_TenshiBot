package models

import (
	"time"
)

// UserMoney représente les informations sur la monnaie d'un utilisateur
type UserMoney struct {
	UserID          string    `json:"user_id"`
	Money           int       `json:"money"`
	LastDailyReward time.Time `json:"last_daily_reward"` // Ajout du champ pour la récompense quotidienne
}
