package controllers

import (
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
)

// CreateShopItem crée un nouvel article dans le magasin
func CreateShopItem(name string, price float64, cooldown int) (*models.ShopItem, error) {
	db := database.DB
	item := models.ShopItem{
		Name:     name,
		Price:    price,
		Cooldown: cooldown,
	}
	result := db.Create(&item)
	return &item, result.Error
}

// GetShopItemByID renvoie un article du magasin par son ID
func GetShopItemByID(id uint) (*models.ShopItem, error) {
	db := database.DB
	var item models.ShopItem
	result := db.First(&item, id)
	return &item, result.Error
}

// GetAllShopItems renvoie tous les articles du magasin
func GetAllShopItems() ([]models.ShopItem, error) {
	db := database.DB
	var items []models.ShopItem
	result := db.Find(&items)
	return items, result.Error
}

// UpdateShopItem met à jour les informations d'un article du magasin
func UpdateShopItem(id uint, name string, price float64, cooldown int) (*models.ShopItem, error) {
	db := database.DB
	var item models.ShopItem
	result := db.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	item.Name = name
	item.Price = price
	item.Cooldown = cooldown
	db.Save(&item)
	return &item, nil
}

// DeleteShopItem supprime un article du magasin par son ID
func DeleteShopItem(id uint) error {
	db := database.DB
	result := db.Delete(&models.ShopItem{}, id)
	return result.Error
}
