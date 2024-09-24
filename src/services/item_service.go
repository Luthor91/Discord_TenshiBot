package services

import (
	"errors"
	"log"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
)

// ItemService est un service pour gérer les items des utilisateurs
type ItemService struct {
	userCtrl *controllers.UserController
	itemCtrl *controllers.ItemController
}

// NewItemService crée une nouvelle instance de ItemService
func NewItemService() *ItemService {
	return &ItemService{
		userCtrl: controllers.NewUserController(),
		itemCtrl: controllers.NewItemController(),
	}
}

// AddItem ajoute un item à l'inventaire de l'utilisateur
func (service *ItemService) AddItem(userID string, itemName string, quantity int) error {
	item := models.Item{
		Name:          itemName,
		Quantity:      quantity,
		UserDiscordID: userID,
	}

	// Vérifier si l'item existe déjà pour cet utilisateur
	existingItem, err := service.itemCtrl.GetItem(userID, itemName)
	if err != nil && !errors.Is(err, controllers.ErrItemNotFound) {
		return err
	}

	if existingItem != nil {
		// L'item existe déjà, augmenter la quantité
		existingItem.Quantity += quantity
		err = service.itemCtrl.UpdateItem(existingItem) // Utilisez ici une méthode Update appropriée
		if err != nil {
			return err
		}
		log.Printf("Ajouté %d %s à l'utilisateur %s", quantity, itemName, userID)
		return nil
	}

	// Si l'item n'existe pas, l'ajouter à l'inventaire
	err = service.itemCtrl.AddItem(&item)
	if err != nil {
		return err
	}

	log.Printf("Ajouté %d %s à l'utilisateur %s", quantity, itemName, userID)
	return nil
}

// GetUserItems récupère les items d'un utilisateur
func (service *ItemService) GetUserItems(userID string) ([]models.Item, error) {
	items, err := service.itemCtrl.GetUserItems(userID)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// RemoveItem retire un certain montant d'un item de l'utilisateur
func (service *ItemService) RemoveItem(userID string, itemName string, amount int) error {
	item, err := service.itemCtrl.GetItem(userID, itemName)
	if err != nil {
		return errors.New("item non trouvé")
	}

	if item.Quantity < amount {
		return errors.New("quantité insuffisante")
	}

	item.Quantity -= amount
	if item.Quantity == 0 {
		// Supprimer l'item s'il n'en reste plus
		err = service.itemCtrl.RemoveItem(item)
		if err != nil {
			return err
		}
		log.Printf("Supprimé l'item %s de l'utilisateur %s", itemName, userID)
	} else {
		err = service.itemCtrl.UpdateItem(item)
		if err != nil {
			return err
		}
		log.Printf("Retiré %d de l'item %s de l'utilisateur %s", amount, itemName, userID)
	}

	return nil
}

// HasItem vérifie si l'utilisateur possède une quantité suffisante de l'item
func (service *ItemService) HasItem(userID, itemName string, quantity int) (bool, error) {
	items, err := service.GetUserItems(userID)
	if err != nil {
		return false, err
	}

	itemCount := 0
	for _, item := range items {
		if item.Name == itemName {
			itemCount += item.Quantity
		}
	}

	return itemCount >= quantity, nil
}

// GiveItem transfère un item d'un utilisateur à un autre
func (service *ItemService) GiveItem(fromUserID, toUserID, itemName string, quantity int) error {
	fromUserItems, err := service.GetUserItems(fromUserID)
	if err != nil {
		return err
	}

	itemToTransfer := &models.Item{}
	itemCount := 0
	for _, item := range fromUserItems {
		if item.Name == itemName {
			itemCount += item.Quantity
			itemToTransfer = &item // Corrigé ici pour référencer le bon item
		}
	}

	if itemCount < quantity {
		return errors.New("pas assez d'items pour transférer")
	}

	// Retirer l'item de l'utilisateur qui donne
	itemToTransfer.Quantity -= quantity
	if itemToTransfer.Quantity == 0 {
		err = service.itemCtrl.RemoveItem(itemToTransfer)
		if err != nil {
			return err
		}
		log.Printf("Item %s transféré de %s à %s", itemName, fromUserID, toUserID)
	} else {
		err = service.itemCtrl.UpdateItem(itemToTransfer)
		if err != nil {
			return err
		}
	}

	// Ajouter l'item à l'utilisateur qui reçoit
	toUserItems, err := service.GetUserItems(toUserID)
	if err != nil {
		return err
	}

	itemFound := false
	for _, item := range toUserItems {
		if item.Name == itemName {
			item.Quantity += quantity
			itemFound = true
			break
		}
	}

	if !itemFound {
		newItem := models.Item{Name: itemName, Quantity: quantity, UserDiscordID: toUserID}
		err = service.itemCtrl.AddItem(&newItem)
		if err != nil {
			return err
		}
	} else {
		err = service.itemCtrl.UpdateItem(itemToTransfer) // On devrait plutôt mettre à jour l'item dans toUser
		if err != nil {
			return err
		}
	}

	log.Printf("Transféré %d de l'item %s de %s à %s", quantity, itemName, fromUserID, toUserID)
	return nil
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
	err = applyItemEffects(target, itemName, quantity)
	if err != nil {
		return err
	}

	// Réduire la quantité de l'item pour l'utilisateur
	user.Items[itemIndex].Quantity -= quantity
	if user.Items[itemIndex].Quantity == 0 {
		user.Items = append(user.Items[:itemIndex], user.Items[itemIndex+1:]...)
	}

	// Sauvegarder les modifications
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
func applyItemEffects(user *models.User, itemName string, quantity int) error {
	// Durée du timeout par item
	const timeoutDurationPerItem = 5 * time.Minute

	if itemName == "timeout" {
		// Calculer la durée totale du timeout basée sur la quantité
		totalTimeoutDuration := timeoutDurationPerItem * time.Duration(quantity)

		// Si l'utilisateur a déjà un timeout en cours
		if !user.TimeoutEnd.IsZero() {
			// Vérifier si le timeout est encore actif
			if time.Now().Before(user.TimeoutEnd) {
				// Ajouter la durée du nouveau timeout à la durée restante
				user.TimeoutEnd = user.TimeoutEnd.Add(totalTimeoutDuration)
			} else {
				// Si le timeout est terminé, définir une nouvelle période
				user.TimeoutEnd = time.Now().Add(totalTimeoutDuration)
			}
		} else {
			// Si l'utilisateur n'a pas de timeout en cours, définir une nouvelle période
			user.TimeoutEnd = time.Now().Add(totalTimeoutDuration)
		}

		log.Printf("Timeout appliqué à l'utilisateur %s jusqu'à %s (durée totale: %s)", user.UserDiscordID, user.TimeoutEnd, totalTimeoutDuration)
	}

	return nil
}
