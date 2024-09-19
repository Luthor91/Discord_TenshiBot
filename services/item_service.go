package services

import (
	"errors"
	"log"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
)

// ItemService est un service pour gérer les items des utilisateurs
type ItemService struct {
	userCtrl *controllers.UserController
}

// NewItemService crée une nouvelle instance de ItemService
func NewItemService() *ItemService {
	return &ItemService{
		userCtrl: &controllers.UserController{DB: database.DB},
	}
}

// AddItem ajoute un item à l'inventaire de l'utilisateur
func (service *ItemService) AddItem(userID string, itemName string, quantity int) error {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	// Trouver l'item dans l'inventaire de l'utilisateur
	itemIndex := -1
	for i, item := range user.Items {
		if item.Name == itemName {
			itemIndex = i
			user.Items[i].Quantity += quantity
			break
		}
	}

	// Si l'item n'existe pas, l'ajouter à l'inventaire
	if itemIndex == -1 {
		newItem := models.Item{
			Name:     itemName,
			Quantity: quantity,
		}
		user.Items = append(user.Items, newItem)
	}

	// Mettre à jour l'utilisateur dans la base de données
	err = service.userCtrl.SaveUser(user)
	if err != nil {
		return err
	}

	log.Printf("Ajouté %d %s à l'utilisateur %s", quantity, itemName, userID)
	return nil
}

// HasItem vérifie si l'utilisateur possède suffisamment de l'item spécifié
func (service *ItemService) HasItem(userID string, itemName string, quantity int) (bool, error) {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		return false, err
	}

	// Vérifier si l'utilisateur a l'item et la quantité nécessaire
	for _, item := range user.Items {
		if item.Name == itemName && item.Quantity >= quantity {
			return true, nil
		}
	}

	return false, nil
}

// UseItem applique un item à un autre utilisateur
func (service *ItemService) UseItem(userID string, targetID string, itemName string, quantity int) error {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		return errors.New("utilisateur non trouvé")
	}

	target, err := service.userCtrl.GetUserByDiscordID(targetID)
	if err != nil {
		return errors.New("cible non trouvée")
	}

	// Vérifier si l'utilisateur possède suffisamment de l'item
	itemIndex := -1
	for i, item := range user.Items {
		if item.Name == itemName {
			itemIndex = i
			if item.Quantity < quantity {
				return errors.New("quantité insuffisante")
			}
			break
		}
	}

	if itemIndex == -1 {
		return errors.New("item non trouvé")
	}

	// Appliquer les effets de l'item
	err = applyItemEffects(service.userCtrl, target, itemName, quantity)
	if err != nil {
		return err
	}

	// Réduire la quantité de l'item pour l'utilisateur
	user.Items[itemIndex].Quantity -= quantity
	if user.Items[itemIndex].Quantity == 0 {
		user.Items = append(user.Items[:itemIndex], user.Items[itemIndex+1:]...)
	}

	// Sauvegarder les modifications dans la base de données
	err = service.userCtrl.SaveUser(user)
	if err != nil {
		return err
	}
	err = service.userCtrl.SaveUser(target)
	if err != nil {
		return err
	}

	log.Printf("Item %s utilisé avec succès sur %s", itemName, targetID)
	return nil
}

// applyItemEffects applique les effets d'un item à un utilisateur
func applyItemEffects(userCtrl *controllers.UserController, user *models.User, itemName string, quantity int) error {
	// Durée du timeout par item
	const timeoutDurationPerItem = 5 * time.Minute

	if itemName == "timeout" {
		// Calculer la durée totale du timeout basée sur la quantité
		totalTimeoutDuration := timeoutDurationPerItem * time.Duration(quantity)

		// Si l'utilisateur a déjà un timeout en cours
		if !user.TimeoutEnd.IsZero() {
			// Ajouter la durée du nouveau timeout à la durée restante
			if time.Now().Before(user.TimeoutEnd) {
				user.TimeoutEnd = user.TimeoutEnd.Add(totalTimeoutDuration)
			} else {
				// Si le timeout est terminé, définir une nouvelle période
				user.TimeoutEnd = time.Now().Add(totalTimeoutDuration)
			}
		} else {
			// Si l'utilisateur n'a pas de timeout en cours, définir une nouvelle période
			user.TimeoutEnd = time.Now().Add(totalTimeoutDuration)
		}

		// Log pour vérifier le timeout
		log.Printf("Timeout appliqué à l'utilisateur %s jusqu'à %s (durée totale: %s)", user.UserDiscordID, user.TimeoutEnd, totalTimeoutDuration)

		// Enregistrer les modifications de l'utilisateur dans la base de données
		err := userCtrl.SaveUser(user)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("effet de l'item non défini")
}
