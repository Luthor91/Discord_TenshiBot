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
	Affinity        int
	Money           int
	Experience      int
	LastDailyReward string
	Rank            int
	RankMoney       int
	RankExperience  int
	RankAffinity    int
	Items           []Item `gorm:"foreignKey:UserDiscordID;references:UserDiscordID"` // Référencer la bonne colonne
	TimeoutEnd      time.Time
}
