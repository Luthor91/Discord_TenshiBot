package models

import (
	"time"

	"gorm.io/gorm"
)

// User représente un utilisateur avec des articles associés
type User struct {
	gorm.Model
	UserID          uint   `gorm:"primarykey"`
	UserDiscordID   string `gorm:"uniqueIndex"`
	Username        string
	Affinity        int
	Money           int
	Experience      int
	LastDailyReward string
	Rank            int
	RankMoney       int
	RankExperience  int
	RankAffinity    int
	Items           []Item `gorm:"foreignKey:UserDiscordID;references:ID"` // Définir la clé étrangère et la clé primaire
	TimeoutEnd      time.Time
}
