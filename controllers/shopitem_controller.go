package controllers

import (
	"errors"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// ShopItemController est un contrôleur pour gérer les opérations sur les articles du magasin
type ShopItemController struct {
	DB *gorm.DB
}

// NewShopItemController crée une nouvelle instance de ShopItemController avec une connexion à la base de données
func NewShopItemController() *ShopItemController {
	return &ShopItemController{
		DB: database.DB,
	}
}

// CreateShopItem crée un nouvel article dans le magasin
func (ctrl *ShopItemController) CreateShopItem(name string, price float64, cooldown int) (*models.ShopItem, error) {
	item := models.ShopItem{
		Name:     name,
		Price:    price,
		Cooldown: cooldown,
	}
	result := ctrl.DB.Create(&item)
	return &item, result.Error
}

// GetShopItemByID renvoie un article du magasin par son ID
func (ctrl *ShopItemController) GetShopItemByID(id uint) (*models.ShopItem, error) {
	var item models.ShopItem
	result := ctrl.DB.First(&item, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Pas de record trouvé, retourner nil sans erreur
	}
	return &item, result.Error // Retourner l'article et l'erreur s'il y en a une
}

// GetAllShopItems renvoie tous les articles du magasin
func (ctrl *ShopItemController) GetAllShopItems() ([]models.ShopItem, error) {
	var items []models.ShopItem
	result := ctrl.DB.Find(&items)
	return items, result.Error
}

// UpdateShopItem met à jour les informations d'un article du magasin
func (ctrl *ShopItemController) UpdateShopItem(id uint, name string, price float64, cooldown int) (*models.ShopItem, error) {
	var item models.ShopItem
	result := ctrl.DB.First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	item.Name = name
	item.Price = price
	item.Cooldown = cooldown
	ctrl.DB.Save(&item)
	return &item, nil
}

// DeleteShopItem supprime un article du magasin par son ID
func (ctrl *ShopItemController) DeleteShopItem(id uint) error {
	result := ctrl.DB.Delete(&models.ShopItem{}, id)
	return result.Error
}
