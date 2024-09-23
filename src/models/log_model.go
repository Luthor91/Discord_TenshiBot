package models

import (
	"time"

	"gorm.io/gorm"
)

// Log représente une entrée de journal
type Log struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	Timestamp     time.Time
	ServerID      string
	ServerName    string
	ChannelID     string
	ChannelName   string
	UserDiscordID string
	Username      string
	Message       string
}
