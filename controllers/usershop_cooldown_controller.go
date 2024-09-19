package controllers

import (
	"time"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
)

// CreateUserShopCooldown crée un nouvel enregistrement de cooldown d'achat pour un utilisateur
func CreateUserShopCooldown(userID string, itemID uint, nextPurchase time.Time) (*models.UserShopCooldown, error) {
	db := database.DB
	user, _ := GetUserIDByDiscordID(userID)
	cooldown := models.UserShopCooldown{
		UserDiscordID: user,
		ItemID:        itemID,
		NextPurchase:  nextPurchase,
	}
	result := db.Create(&cooldown)
	return &cooldown, result.Error
}

// GetUserShopCooldown renvoie le cooldown d'achat d'un utilisateur pour un article spécifique
func GetUserShopCooldown(userID string, itemID uint) (*models.UserShopCooldown, error) {
	db := database.DB
	var cooldown models.UserShopCooldown
	result := db.First(&cooldown, "user_id = ? AND item_id = ?", userID, itemID)
	return &cooldown, result.Error
}

// UpdateUserShopCooldown met à jour le délai d'achat pour un utilisateur pour un article spécifique
func UpdateUserShopCooldown(userID string, itemID uint, nextPurchase time.Time) (*models.UserShopCooldown, error) {
	db := database.DB
	var cooldown models.UserShopCooldown
	result := db.First(&cooldown, "user_id = ? AND item_id = ?", userID, itemID)
	if result.Error != nil {
		return nil, result.Error
	}
	cooldown.NextPurchase = nextPurchase
	db.Save(&cooldown)
	return &cooldown, nil
}

// DeleteUserShopCooldown supprime le délai d'achat d'un utilisateur pour un article spécifique
func DeleteUserShopCooldown(userID string, itemID uint) error {
	db := database.DB
	result := db.Delete(&models.UserShopCooldown{}, "user_id = ? AND item_id = ?", userID, itemID)
	return result.Error
}
