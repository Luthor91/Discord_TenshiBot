package controllers

import (
	"errors"

	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/models"
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
	return result.Error == nil, result.Error // Retourne vrai si l'utilisateur existe
}

// UserExistsByDiscordID vérifie si un utilisateur existe en utilisant son ID Discord
func (controller *UserController) UserExistsByDiscordID(userDiscordID string) (bool, error) {
	var user models.User
	result := controller.DB.Where("user_discord_id = ?", userDiscordID).Find(&user)

	// Vérifie si des lignes ont été trouvées
	if result.RowsAffected == 0 {
		return false, nil // L'utilisateur n'existe pas
	}

	// Retourne vrai si l'utilisateur existe, ou une erreur si elle s'est produite
	return true, result.Error
}

// GetAllUsers récupère tous les utilisateurs
func (ctrl *UserController) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := ctrl.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserIDByDiscordID récupère l'identifiant interne de l'utilisateur en utilisant son ID Discord
func (ctrl *UserController) GetUserIDByDiscordID(discordID string) (uint, error) {
	var user models.User
	if err := ctrl.DB.First(&user, "user_discord_id = ?", discordID).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetUserDiscordIDByID récupère l'ID Discord de l'utilisateur en utilisant son identifiant interne
func (ctrl *UserController) GetUserDiscordIDByID(userID uint) (string, error) {
	var user models.User
	if err := ctrl.DB.First(&user, userID).Error; err != nil {
		return "", err
	}
	return user.UserDiscordID, nil
}

// UpdateUser met à jour les informations d'un utilisateur dans la base de données
func (ctrl *UserController) UpdateUser(user *models.User) error {
	return ctrl.DB.Save(user).Error
}

// CreateUser crée un nouvel utilisateur
func (ctrl *UserController) CreateUser(userID, username string) (*models.User, error) {
	user := models.User{
		UserDiscordID:   userID,
		Username:        username,
	}
	// Utiliser FirstOrCreate pour vérifier l'existence de l'utilisateur
	if err := ctrl.DB.Where(models.User{UserDiscordID: userID}).FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByDiscordID récupère un utilisateur par son ID Discord
func (ctrl *UserController) GetUserByDiscordID(userDiscordID string) (*models.User, error) {
	var user models.User
	result := ctrl.DB.Where("user_discord_id = ?", userDiscordID).Find(&user)

	// Vérifie si l'utilisateur existe
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Retourne l'utilisateur si trouvé ou l'erreur s'il y en a une
	return &user, result.Error
}

// GetUserByID récupère un utilisateur par son identifiant interne
func (ctrl *UserController) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	result := ctrl.DB.Find(&user, userID)

	// Vérifie si l'utilisateur existe
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Retourne l'utilisateur si trouvé ou l'erreur s'il y en a une
	return &user, result.Error
}

// SaveUser met à jour ou insère un utilisateur dans la base de données
func (ctrl *UserController) SaveUser(user *models.User) error {
	var existingUser models.User
	if err := ctrl.DB.First(&existingUser, "user_discord_id = ?", user.UserDiscordID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Si l'utilisateur n'existe pas, créer un nouvel enregistrement
			return ctrl.DB.Create(user).Error
		}
		return err
	}

	// Mettre à jour l'utilisateur existant
	existingUser.Username = user.Username

	// Sauvegarder les modifications
	return ctrl.DB.Save(&existingUser).Error
}

// DeleteUser supprime un utilisateur
func (ctrl *UserController) DeleteUser(userID string) error {
	return ctrl.DB.Delete(&models.User{}, "user_discord_id = ?", userID).Error
}


// AddUserIfNotExists ajoute un utilisateur s'il n'existe pas déjà
func (ctrl *UserController) AddUserIfNotExists(userID, username string) error {
	exists, err := ctrl.UserExistsByDiscordID(userID)
	if err != nil {
		return err
	}
	if !exists {
		_, err := ctrl.CreateUser(userID, username)
		return err
	}
	return nil
}
