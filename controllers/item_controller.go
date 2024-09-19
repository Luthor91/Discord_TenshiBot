package controllers

import (
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// ItemController est un contrôleur pour gérer les opérations sur les items
type ItemController struct {
	DB *gorm.DB
}

// CreateItem crée un nouvel item
func (ctrl *ItemController) CreateItem(name string, quantity int) (*models.Item, error) {
	item := models.Item{Name: name, Quantity: quantity}
	if err := ctrl.DB.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetItem récupère un item par ID
func (ctrl *ItemController) GetItem(id uint) (*models.Item, error) {
	var item models.Item
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateItem met à jour un item
func (ctrl *ItemController) UpdateItem(id uint, name string, quantity int) (*models.Item, error) {
	var item models.Item
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		return nil, err
	}
	item.Name = name
	item.Quantity = quantity
	if err := ctrl.DB.Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// DeleteItem supprime un item
func (ctrl *ItemController) DeleteItem(id uint) error {
	if err := ctrl.DB.Delete(&models.Item{}, id).Error; err != nil {
		return err
	}
	return nil
}
