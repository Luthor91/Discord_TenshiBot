package models

import (
	"gorm.io/gorm"
)

type ReactionRole struct {
	gorm.Model
	GuildID   string `gorm:"index;not null"`
	ChannelID string `gorm:"not null"`
	MessageID string `gorm:"index;not null"`
	Emoji     string `gorm:"not null"`
	RoleID    string `gorm:"not null"`
}
