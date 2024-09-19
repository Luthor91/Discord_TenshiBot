package controllers

import (
	"errors"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// UserController est un contrôleur pour gérer les opérations sur les utilisateurs
type UserController struct {
	DB *gorm.DB
}

// NewUserController crée une nouvelle instance de UserController avec une connexion à la base de données
func NewUserController() *UserController {
	return &UserController{
		DB: database.DB,
	}
}

// UserExistsByID vérifie si un utilisateur existe en utilisant son ID
func (controller *UserController) UserExistsByID(userID uint) (bool, error) {
	var user models.User
	result := controller.DB.First(&user, userID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil // L'utilisateur n'existe pas
	}
	return true, result.Error // L'utilisateur existe
}

func (controller *UserController) UserExistsByDiscordID(userDiscordID string) (bool, error) {
	var user models.User
	result := controller.DB.Where("user_discord_id = ?", userDiscordID).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil // L'utilisateur n'existe pas
	}
	return true, result.Error // L'utilisateur existe
}

// GetUserIDByDiscordID récupère l'identifiant interne de l'utilisateur en utilisant son ID Discord
func (ctrl *UserController) GetUserIDByDiscordID(discordID string) (uint, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "user_id = ?", discordID).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserDiscordIDByID récupère l'ID Discord de l'utilisateur en utilisant son identifiant interne
func (ctrl *UserController) GetUserDiscordIDByID(userID uint) (string, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "id = ?", userID).Error; err != nil {
		return "", err
	}
	return user.UserDiscordID, nil
}

// UpdateUser met à jour les informations d'un utilisateur dans la base de données
func (ctrl *UserController) UpdateUser(updatedUser models.User) bool {
	result := ctrl.DB.Model(&models.User{}).Where("user_id = ?", updatedUser.UserDiscordID).Updates(updatedUser)
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

// GetUserByDiscordID récupère un utilisateur par son ID Discord
func (ctrl *UserController) GetUserByDiscordID(userDiscordID string) (*models.User, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "user_discord_id = ?", userDiscordID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID récupère un utilisateur par son identifiant interne
func (ctrl *UserController) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveUser met à jour ou insère un utilisateur dans la base de données
func (ctrl *UserController) SaveUser(user *models.User) error {
	var existingUser models.User
	if err := ctrl.DB.First(&existingUser, "user_id = ?", user.UserDiscordID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, créer un nouvel enregistrement
			if err := ctrl.DB.Create(user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		// Mettre à jour l'utilisateur existant
		existingUser.Username = user.Username
		existingUser.Affinity = user.Affinity
		existingUser.Money = user.Money
		existingUser.Experience = user.Experience
		existingUser.LastDailyReward = user.LastDailyReward
		existingUser.Rank = user.Rank
		existingUser.RankMoney = user.RankMoney
		existingUser.RankExperience = user.RankExperience
		existingUser.RankAffinity = user.RankAffinity

		// Sauvegarder les modifications
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
