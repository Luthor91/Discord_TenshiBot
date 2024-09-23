package models

import "gorm.io/gorm"

// ShopItem représente un article dans le magasin
type ShopItem struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey"`
	Name     string  `gorm:"unique;not null"`
	Price    float64 `gorm:"not null"`
	Emoji    string  `gorm:"unique;not null"`
	Cooldown int     `gorm:"not null"`
}
