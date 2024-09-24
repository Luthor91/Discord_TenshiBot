package controllers

import (
	"errors"
	"time"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// ShopController est un contrôleur pour gérer les opérations sur les articles du magasin
type ShopController struct {
	DB *gorm.DB
}

// NewShopController crée une nouvelle instance de ShopController avec une connexion à la base de données
func NewShopController() *ShopController {
	return &ShopController{
		DB: database.DB,
	}
}

// CreateShopItem crée un nouvel article dans le magasin
func (ctrl *ShopController) CreateShopItem(name string, price float64, cooldown int) (*models.ShopItem, error) {
	item := models.ShopItem{
		Name:     name,
		Price:    price,
		Cooldown: cooldown,
	}
	result := ctrl.DB.Create(&item)
	if result.Error != nil {
		return nil, result.Error // Retourner une erreur si la création échoue
	}
	return &item, nil
}

// GetShopItemByID renvoie un article du magasin par son ID
func (ctrl *ShopController) GetShopItemByID(id uint) (*models.ShopItem, error) {
	var item models.ShopItem
	result := ctrl.DB.First(&item, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // Pas de record trouvé, retourner nil sans erreur
	}
	return &item, result.Error // Retourner l'article et l'erreur s'il y en a une
}

// GetShopItemByName récupère un item de la boutique par son nom
func (ctrl *ShopController) GetShopItemByName(name string) (*models.ShopItem, error) {
	var item models.ShopItem
	result := ctrl.DB.Where("name = ? AND deleted_at IS NULL", name).First(&item)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("article non trouvé")
	}

	return &item, result.Error
}

// GetAllShopItems renvoie tous les articles du magasin
func (ctrl *ShopController) GetAllShopItems() ([]models.ShopItem, error) {
	var items []models.ShopItem
	result := ctrl.DB.Find(&items)
	return items, result.Error
}

// UpdateShopItem met à jour les informations d'un article du magasin
func (ctrl *ShopController) UpdateShopItem(id uint, name string, price float64, cooldown int) (*models.ShopItem, error) {
	var item models.ShopItem
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		return nil, err // Retourner l'erreur si l'article n'est pas trouvé
	}

	item.Name = name
	item.Price = price
	item.Cooldown = cooldown
	if err := ctrl.DB.Save(&item).Error; err != nil {
		return nil, err // Retourner une erreur si la mise à jour échoue
	}
	return &item, nil
}

// DeleteShopItem supprime un article du magasin par son ID
func (ctrl *ShopController) DeleteShopItem(id uint) error {
	result := ctrl.DB.Delete(&models.ShopItem{}, id)
	return result.Error
}

// CreateUserShopCooldown crée un nouvel enregistrement de cooldown d'achat pour un utilisateur
func (ctrl *ShopController) CreateUserShopCooldown(userCooldown *models.UserShopCooldown) error {
	if err := ctrl.DB.Create(userCooldown).Error; err != nil {
		return err
	}
	return nil
}

// GetUserShopCooldown renvoie le cooldown d'achat d'un utilisateur pour un article spécifique
func (ctrl *ShopController) GetUserShopCooldown(userDiscordID string, itemID uint) (*models.UserShopCooldown, error) {
	var cooldown models.UserShopCooldown
	result := ctrl.DB.First(&cooldown, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Si l'enregistrement n'existe pas, crée un nouvel enregistrement par défaut
		cooldown = models.UserShopCooldown{
			UserDiscordID: userDiscordID, // Assurez-vous que cette ligne est correcte
			ItemID:        itemID,
			NextPurchase:  time.Time{}, // Définit un NextPurchase par défaut
		}
		if err := ctrl.DB.Create(&cooldown).Error; err != nil {
			return nil, err
		}
		return &cooldown, nil // Retourne le nouveau cooldown créé
	}

	return &cooldown, result.Error
}

// SetUserShopCooldown définit ou met à jour le cooldown d'achat pour un utilisateur
func (controller *ShopController) SetUserShopCooldown(userID string, itemID uint, nextPurchase time.Time) error {
	cooldown := &models.UserShopCooldown{
		UserDiscordID: userID,
		ItemID:        itemID,
		NextPurchase:  nextPurchase,
	}
	result := controller.DB.Create(cooldown)
	return result.Error
}

// UpdateUserShopCooldown met à jour le délai d'achat pour un utilisateur pour un article spécifique
func (ctrl *ShopController) UpdateUserShopCooldown(userDiscordID string, itemID uint, nextPurchase time.Time) error {
	// Chercher l'enregistrement du cooldown
	var cooldown models.UserShopCooldown
	if err := ctrl.DB.First(&cooldown, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID).Error; err != nil {
		return err
	}

	// Mettre à jour le cooldown et enregistrer les changements
	cooldown.NextPurchase = nextPurchase
	return ctrl.DB.Save(&cooldown).Error
}

// DeleteUserShopCooldown supprime le délai d'achat d'un utilisateur pour un article spécifique
func (ctrl *ShopController) DeleteUserShopCooldown(userDiscordID string, itemID uint) error {
	result := ctrl.DB.Delete(&models.UserShopCooldown{}, "user_discord_id = ? AND item_id = ?", userDiscordID, itemID)
	return result.Error
}
