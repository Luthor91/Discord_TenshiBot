package services

import (
	"time"

	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// AddMoney ajoute de la monnaie à un utilisateur donné dans la base de données
func (service *UserService) AddMoney(userID string, amount int) error {
	var user models.User

	// Verrouillage pour éviter les accès concurrents
	service.mu.Lock()
	defer service.mu.Unlock()

	// Rechercher l'utilisateur dans la base de données
	if err := service.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, le créer
			user = models.User{
				UserDiscordID: userID,
				Money:         amount,
				Affinity:      0,
				Experience:    0,
			}
			if err := service.db.Create(&user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Si l'utilisateur existe, ajouter de la monnaie
		user.Money += amount
		if err := service.db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetUserMoney renvoie la quantité de monnaie d'un utilisateur depuis la base de données
func (service *UserService) GetUserMoney(userID string) (int, error) {
	var user models.User

	// Rechercher l'utilisateur dans la base de données
	if err := service.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, retourner 0
			return 0, nil
		}
		return 0, err
	}

	return user.Money, nil
}

// CanReceiveDailyReward vérifie si l'utilisateur peut recevoir une récompense quotidienne
func (service *UserService) CanReceiveDailyReward(userID string) (bool, time.Duration, error) {
	var user models.User

	// Rechercher l'utilisateur dans la base de données
	if err := service.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, il peut recevoir une récompense
			return true, 0, nil
		}
		return false, 0, err
	}

	// Convertir la chaîne de la dernière récompense en `time.Time`
	lastRewardTime, err := time.Parse(time.RFC3339, user.LastDailyReward)
	if err != nil {
		// Si le parsing échoue, assumer que l'utilisateur peut recevoir une récompense
		return true, 0, nil
	}

	now := time.Now()
	if now.Sub(lastRewardTime).Hours() >= 24 {
		return true, 0, nil
	}
	return false, time.Until(lastRewardTime.Add(24 * time.Hour)), nil
}

// GiveDailyMoney accorde la récompense quotidienne et met à jour la date de la dernière récompense
func (service *UserService) GiveDailyMoney(userID string, amount int) error {
	var user models.User

	// Verrouiller le mutex pour empêcher les accès concurrents
	service.mu.Lock()
	defer service.mu.Unlock()

	// Rechercher l'utilisateur dans la base de données
	if err := service.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, le créer
			user = models.User{
				UserDiscordID:   userID,
				Money:           amount,
				Affinity:        0,
				Experience:      0,
				LastDailyReward: time.Now().Format(time.RFC3339),
			}
			if err := service.db.Create(&user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Si l'utilisateur existe, ajouter de la monnaie et mettre à jour la date de la dernière récompense
		user.Money += amount
		user.LastDailyReward = time.Now().Format(time.RFC3339)
		if err := service.db.Save(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
