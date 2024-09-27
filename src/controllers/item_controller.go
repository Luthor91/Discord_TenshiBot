package controllers

import (
	"errors"
	"log"

	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/models"
	"gorm.io/gorm"
)

// ItemController gère les opérations sur les items
type ItemController struct {
	DB *gorm.DB
}

// NewItemController crée une nouvelle instance de ItemController
func NewItemController() *ItemController {
	return &ItemController{DB: database.DB} // Assurez-vous que db est votre instance gorm.DB
}

// AddItem ajoute un nouvel item à la base de données
func (ctrl *ItemController) AddItem(item *models.Item) error {
	if err := ctrl.DB.Create(item).Error; err != nil {
		log.Printf("Erreur lors de l'ajout de l'item %s : %v", item.Name, err)
		return err
	}
	return nil
}

// GetItem récupère un item par son nom et l'ID de l'utilisateur
func (ctrl *ItemController) GetItem(userID string, itemName string) (*models.Item, error) {
	var item models.Item
	if err := ctrl.DB.Where("user_discord_id = ? AND name = ?", userID, itemName).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		log.Printf("Erreur lors de la récupération de l'item %s pour l'utilisateur %s : %v", itemName, userID, err)
		return nil, err
	}
	return &item, nil
}

// GetUserItems récupère tous les items d'un utilisateur
func (ctrl *ItemController) GetUserItems(userID string) ([]models.Item, error) {
	var items []models.Item
	if err := ctrl.DB.Where("user_discord_id = ?", userID).Find(&items).Error; err != nil {
		log.Printf("Erreur lors de la récupération des items pour l'utilisateur %s : %v", userID, err)
		return nil, err
	}
	return items, nil
}

// UpdateItem met à jour un item existant dans la base de données
func (ctrl *ItemController) UpdateItem(item *models.Item) error {
	if err := ctrl.DB.Save(item).Error; err != nil {
		log.Printf("Erreur lors de la mise à jour de l'item %s : %v", item.Name, err)
		return err
	}
	return nil
}

// RemoveItem supprime un item de la base de données
func (ctrl *ItemController) RemoveItem(item *models.Item) error {
	if err := ctrl.DB.Delete(item).Error; err != nil {
		log.Printf("Erreur lors de la suppression de l'item %s : %v", item.Name, err)
		return err
	}
	return nil
}

// ErrItemNotFound est l'erreur retournée lorsque l'item n'est pas trouvé
var ErrItemNotFound = errors.New("item non trouvé")
