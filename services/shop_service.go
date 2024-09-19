package services

import (
	"errors"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// ShopService est un service pour gérer les opérations liées aux achats en boutique
type ShopService struct {
	db *gorm.DB
}

// NewShopService crée une nouvelle instance de ShopService
func NewShopService() *ShopService {
	return &ShopService{
		db: database.DB,
	}
}

// LoadUserShopCooldowns récupère tous les cooldowns d'achats pour un utilisateur spécifique
func (service *ShopService) LoadUserShopCooldowns(userDiscordID string) (map[uint]models.UserShopCooldown, error) {
	var userCooldowns []models.UserShopCooldown

	// Chercher tous les cooldowns pour cet utilisateur
	result := service.db.Where("user_discord_id = ?", userDiscordID).Find(&userCooldowns)
	if result.Error != nil {
		return nil, result.Error
	}

	// Créer un mappage des cooldowns par ID d'article
	cooldownMap := make(map[uint]models.UserShopCooldown)
	for _, cooldown := range userCooldowns {
		cooldownMap[cooldown.ItemID] = cooldown
	}

	return cooldownMap, nil
}

// SetShopCooldown définit le cooldown pour un utilisateur sur un article spécifique
func (service *ShopService) SetShopCooldown(userDiscordID string, itemID uint, nextPurchase time.Time) error {
	// On vérifie si l'article existe dans la base de données
	item, err := controllers.NewShopItemController().GetShopItemByID(itemID)
	if err != nil {
		return err
	}

	userCooldown, err := controllers.NewUserShopCooldownController().GetUserShopCooldown(userDiscordID, itemID)
	if err != nil {
		return err
	}

	if userCooldown.UserDiscordID == "" {
		// Créer un nouvel enregistrement de cooldown si l'utilisateur n'en a pas encore
		exists, _ := controllers.NewUserController().UserExistsByDiscordID(userDiscordID)
		if !exists {
			return errors.New("utilisateur non trouvé")
		}
		userCooldown = &models.UserShopCooldown{
			UserDiscordID: userCooldown.UserDiscordID,
			ItemID:        item.ID,
			NextPurchase:  nextPurchase,
		}
		if err := service.db.Create(&userCooldown).Error; err != nil {
			return err
		}
	} else {
		// Mettre à jour le cooldown si l'utilisateur a déjà un enregistrement
		userCooldown.NextPurchase = nextPurchase
		if err := service.db.Save(&userCooldown).Error; err != nil {
			return err
		}
	}

	return nil
}

// IsCooldownExpired vérifie si le cooldown est expiré pour un utilisateur et un article spécifique
func (service *ShopService) IsCooldownExpired(userID string, itemID uint) (bool, error) {
	// Chercher l'enregistrement du cooldown pour cet utilisateur et cet article
	var userCooldown models.UserShopCooldown
	result := service.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&userCooldown)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Aucun cooldown trouvé, l'utilisateur peut acheter l'article
		return true, nil
	}

	// Vérifier si le cooldown est expiré
	return time.Now().After(userCooldown.NextPurchase), nil
}

// GetShopItems récupère tous les items de la boutique
func (service *ShopService) GetShopItems() ([]models.ShopItem, error) {
	var items []models.ShopItem
	if err := service.db.Where("deleted_at IS NULL").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetShopCooldown renvoie le temps restant avant que l'utilisateur puisse acheter à nouveau un article
func (service *ShopService) GetShopCooldown(userID string, itemID uint) (time.Time, error) {
	// Chercher l'enregistrement du cooldown pour cet utilisateur et cet article
	var userCooldown models.UserShopCooldown
	result := service.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&userCooldown)

	if result.Error != nil {
		return time.Time{}, result.Error
	}

	// Retourner le prochain moment où l'utilisateur pourra acheter l'article
	return userCooldown.NextPurchase, nil
}
