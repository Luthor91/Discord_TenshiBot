package controllers

import (
	"errors"
	"time"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// UserShopCooldownController gère les opérations sur les délais d'achat des utilisateurs
type UserShopCooldownController struct {
	DB *gorm.DB
}

// NewUserShopCooldownController crée une nouvelle instance de UserShopCooldownController
func NewUserShopCooldownController() *UserShopCooldownController {
	return &UserShopCooldownController{
		DB: database.DB,
	}
}

// CreateUserShopCooldown crée un nouvel enregistrement de cooldown d'achat pour un utilisateur
func (ctrl *UserShopCooldownController) CreateUserShopCooldown(userID string, itemID uint, nextPurchase time.Time) (*models.UserShopCooldown, error) {
	_, err := NewUserController().GetUserIDByDiscordID(userID)
	if err != nil {
		return nil, err
	}
	cooldown := models.UserShopCooldown{
		UserDiscordID: userID,
		ItemID:        itemID,
		NextPurchase:  nextPurchase,
	}
	result := ctrl.DB.Create(&cooldown)
	return &cooldown, result.Error
}

// GetUserShopCooldown renvoie le cooldown d'achat d'un utilisateur pour un article spécifique
func (ctrl *UserShopCooldownController) GetUserShopCooldown(userDiscordID string, itemID uint) (*models.UserShopCooldown, error) {
	var cooldown models.UserShopCooldown
	result := ctrl.DB.First(&cooldown, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Si l'enregistrement n'existe pas, crée un nouvel enregistrement par défaut
		cooldown = models.UserShopCooldown{
			UserDiscordID: userDiscordID,
			ItemID:        itemID,
			NextPurchase:  time.Time{}, // Définit un NextPurchase par défaut
		}
		if err := ctrl.DB.Create(&cooldown).Error; err != nil {
			return nil, err
		}
	}

	return &cooldown, result.Error
}

// SetUserShopCooldown définit ou met à jour le cooldown d'achat pour un utilisateur
func (ctrl *UserShopCooldownController) SetUserShopCooldown(userDiscordID string, itemID uint, nextPurchase time.Time) error {
	var cooldown models.UserShopCooldown
	result := ctrl.DB.First(&cooldown, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Créer un nouvel enregistrement si aucun cooldown n'existe
		cooldown = models.UserShopCooldown{
			UserDiscordID: userDiscordID,
			ItemID:        itemID,
			NextPurchase:  nextPurchase,
		}
		return ctrl.DB.Create(&cooldown).Error
	} else if result.Error != nil {
		// Erreur lors de la récupération du cooldown
		return result.Error
	}

	// Sinon, mettre à jour l'enregistrement existant
	cooldown.NextPurchase = nextPurchase
	return ctrl.DB.Save(&cooldown).Error
}

// UpdateUserShopCooldown met à jour le délai d'achat pour un utilisateur pour un article spécifique
func (ctrl *UserShopCooldownController) UpdateUserShopCooldown(userDiscordID string, itemID uint, nextPurchase time.Time) (*models.UserShopCooldown, error) {
	var cooldown models.UserShopCooldown
	result := ctrl.DB.First(&cooldown, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)
	if result.Error != nil {
		return nil, result.Error
	}
	cooldown.NextPurchase = nextPurchase
	if err := ctrl.DB.Save(&cooldown).Error; err != nil {
		return nil, err
	}
	return &cooldown, nil
}

// DeleteUserShopCooldown supprime le délai d'achat d'un utilisateur pour un article spécifique
func (ctrl *UserShopCooldownController) DeleteUserShopCooldown(userDiscordID string, itemID uint) error {
	result := ctrl.DB.Delete(&models.UserShopCooldown{}, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)
	return result.Error
}
