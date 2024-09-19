package models

import (
	"gorm.io/gorm"
)

// BadWord représente un mot considéré comme négatif
type BadWord struct {
	gorm.Model
	Word string `gorm:"unique;not null"`
}

// GoodWord représente un mot considéré comme positif
type GoodWord struct {
	gorm.Model
	Word string `gorm:"unique;not null"`
}
