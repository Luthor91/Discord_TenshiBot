package controllers

import (
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// UserController est un contrôleur pour gérer les opérations sur les utilisateurs
type UserController struct {
	DB *gorm.DB
}

// GetUserIDByDiscordID récupère l'identifiant interne de l'utilisateur en utilisant son ID Discord
func GetUserIDByDiscordID(discordID string) (uint, error) {
	var user models.User
	if err := database.DB.First(&user, "user_id = ?", discordID).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserDiscordIDByID récupère l'ID Discord de l'utilisateur en utilisant son identifiant interne
func GetUserDiscordIDByID(userID uint) (string, error) {
	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		return "", err
	}
	return user.UserDiscordID, nil
}

// UpdateUser met à jour les informations d'un utilisateur dans la base de données
func UpdateUser(userID string, updatedUser models.User) bool {
	result := database.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updatedUser)
	return result.RowsAffected > 0
}

// CreateUser crée un nouvel utilisateur
func (ctrl *UserController) CreateUser(userID, username string, affinity, money, experience int, lastDailyReward string, rank, rankMoney, rankExperience, rankAffinity int) (*models.User, error) {
	user := models.User{
		UserDiscordID:   userID,
		Username:        username,
		Affinity:        affinity,
		Money:           money,
		Experience:      experience,
		LastDailyReward: lastDailyReward,
		Rank:            rank,
		RankMoney:       rankMoney,
		RankExperience:  rankExperience,
		RankAffinity:    rankAffinity,
	}
	if err := ctrl.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser récupère un utilisateur par ID
func (ctrl *UserController) GetUserByDiscordID(userID string) (*models.User, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser récupère un utilisateur par ID
func (ctrl *UserController) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser met à jour un utilisateur
func (ctrl *UserController) UpdateUser(userID string, username string, affinity, money, experience int, lastDailyReward string, rank, rankMoney, rankExperience, rankAffinity int) (*models.User, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	user.Username = username
	user.Affinity = affinity
	user.Money = money
	user.Experience = experience
	user.LastDailyReward = lastDailyReward
	user.Rank = rank
	user.RankMoney = rankMoney
	user.RankExperience = rankExperience
	user.RankAffinity = rankAffinity
	if err := ctrl.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveUser met à jour ou insère un utilisateur dans la base de données
func (ctrl *UserController) SaveUser(user *models.User) error {
	// Cherche si l'utilisateur existe déjà dans la base de données
	var existingUser models.User
	if err := ctrl.DB.First(&existingUser, "user_id = ?", user.UserDiscordID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, créer un nouvel enregistrement
			if err := ctrl.DB.Create(user).Error; err != nil {
				return err
			}
		} else {
			// Autre erreur lors de la recherche de l'utilisateur
			return err
		}
	} else {
		// Met à jour les champs de l'utilisateur existant
		existingUser.Username = user.Username
		existingUser.Affinity = user.Affinity
		existingUser.Money = user.Money
		existingUser.Experience = user.Experience
		existingUser.LastDailyReward = user.LastDailyReward
		existingUser.Rank = user.Rank
		existingUser.RankMoney = user.RankMoney
		existingUser.RankExperience = user.RankExperience
		existingUser.RankAffinity = user.RankAffinity

		// Sauvegarde les modifications
		if err := ctrl.DB.Save(&existingUser).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteUser supprime un utilisateur
func (ctrl *UserController) DeleteUser(userID string) error {
	if err := ctrl.DB.Delete(&models.User{}, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
