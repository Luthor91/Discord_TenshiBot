package services

import (
	"errors"
	"log"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
)

// ShopService est un service pour gérer les opérations liées aux achats en boutique
type ShopService struct {
	shopController *controllers.ShopController
}

// NewShopService crée une nouvelle instance de ShopService
func NewShopService() *ShopService {
	shopController := controllers.NewShopController()
	return &ShopService{
		shopController: shopController,
	}
}

// LoadUserShopCooldowns récupère tous les cooldowns d'achats pour un utilisateur spécifique
func (service *ShopService) LoadUserShopCooldowns(userDiscordID string) (map[uint]models.UserShopCooldown, error) {
	var userCooldowns []models.UserShopCooldown

	// Chercher tous les cooldowns pour cet utilisateur
	result := service.shopController.DB.Where("user_discord_id = ?", userDiscordID).Find(&userCooldowns)
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
	// Vérifier si l'article existe dans la base de données
	item, err := service.shopController.GetShopItemByID(itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("article non trouvé")
	}

	userCooldown, err := service.shopController.GetUserShopCooldown(userDiscordID, itemID)
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
			UserDiscordID: userDiscordID,
			ItemID:        item.ID,
			NextPurchase:  nextPurchase,
		}

		// Utiliser la méthode CreateUserShopCooldown pour l'enregistrer
		if err := service.shopController.CreateUserShopCooldown(userCooldown); err != nil {
			return err
		}
	} else {
		// Mettre à jour le cooldown si l'utilisateur a déjà un enregistrement
		userCooldown.NextPurchase = nextPurchase
		if err := service.shopController.UpdateUserShopCooldown(userDiscordID, itemID, nextPurchase); err != nil {
			return err
		}
	}

	return nil
}

// IsCooldownExpired vérifie si le cooldown est expiré pour un utilisateur et un article spécifique
func (service *ShopService) IsCooldownExpired(userDiscordID string, itemID uint) (bool, error) {
	// Chercher l'enregistrement du cooldown pour cet utilisateur et cet article
	userCooldown, err := service.shopController.GetUserShopCooldown(userDiscordID, itemID)
	if err != nil {
		return false, err
	}

	if userCooldown == nil {
		// Aucun cooldown trouvé, l'utilisateur peut acheter l'article
		return true, nil
	}

	// Vérifier si le cooldown est expiré
	return time.Now().After(userCooldown.NextPurchase), nil
}

// GetShopItems récupère tous les items de la boutique
func (service *ShopService) GetShopItems() ([]models.ShopItem, error) {
	return service.shopController.GetAllShopItems()
}

// GetShopItemByName récupère un item de la boutique par son nom
func (service *ShopService) GetShopItemByName(name string) (*models.ShopItem, error) {
	return service.shopController.GetShopItemByName(name)
}

// GetShopCooldown renvoie le temps restant avant que l'utilisateur puisse acheter à nouveau un article
func (service *ShopService) GetShopCooldown(userDiscordID string, itemID uint) (time.Time, error) {
	userCooldown, err := service.shopController.GetUserShopCooldown(userDiscordID, itemID)
	if err != nil {
		return time.Time{}, err
	}

	if userCooldown == nil {
		// Aucun cooldown trouvé
		return time.Time{}, errors.New("aucun cooldown trouvé")
	}

	// Retourner le prochain moment où l'utilisateur pourra acheter l'article
	return userCooldown.NextPurchase, nil
}

// GetUserShopCooldown récupère le cooldown d'achat pour un utilisateur donné et un item spécifique.
func (s *ShopService) GetUserShopCooldown(userID string, itemID uint) (*models.UserShopCooldown, error) {
	cooldown, err := controllers.NewShopController().GetUserShopCooldown(userID, itemID)
	if err != nil {
		log.Println("Erreur dans le service GetUserShopCooldown:", err)
		return nil, err
	}
	return cooldown, nil
}

// SetUserShopCooldown met à jour le cooldown d'achat pour un utilisateur donné et un item spécifique.
func (s *ShopService) SetUserShopCooldown(userID string, itemID uint, nextPurchaseTime time.Time) error {
	err := controllers.NewShopController().SetUserShopCooldown(userID, itemID, nextPurchaseTime)
	if err != nil {
		log.Println("Erreur dans le service SetUserShopCooldown:", err)
		return err
	}
	return nil
}
