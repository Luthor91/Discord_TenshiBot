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
	existingUser.Affinity = user.Affinity
	existingUser.Money = user.Money
	existingUser.Experience = user.Experience
	existingUser.LastDailyReward = user.LastDailyReward
	existingUser.Rank = user.Rank
	existingUser.RankMoney = user.RankMoney
	existingUser.RankExperience = user.RankExperience
	existingUser.RankAffinity = user.RankAffinity

	// Sauvegarder les modifications
	return ctrl.DB.Save(&existingUser).Error
}

// DeleteUser supprime un utilisateur
func (ctrl *UserController) DeleteUser(userID string) error {
	return ctrl.DB.Delete(&models.User{}, "user_discord_id = ?", userID).Error
}

// GiveMoney transfère une somme d'argent d'un utilisateur à un autre
func (ctrl *UserController) GiveMoney(fromUserID, toUserID string, moneyAmount int) error {
	fromUser, err := ctrl.GetUserByDiscordID(fromUserID)
	if err != nil {
		return err
	}

	toUser, err := ctrl.GetUserByDiscordID(toUserID)
	if err != nil {
		return err
	}

	if fromUser.Money < moneyAmount {
		return errors.New("not enough money to transfer")
	}

	fromUser.Money -= moneyAmount
	toUser.Money += moneyAmount

	if err := ctrl.UpdateUser(fromUser); err != nil {
		return err
	}
	if err := ctrl.UpdateUser(toUser); err != nil {
		return err
	}

	return nil
}

// GiveExperience transfère une quantité d'expérience d'un utilisateur à un autre
func (ctrl *UserController) GiveExperience(fromUserID, toUserID string, xpAmount int) error {
	fromUser, err := ctrl.GetUserByDiscordID(fromUserID)
	if err != nil {
		return err
	}

	toUser, err := ctrl.GetUserByDiscordID(toUserID)
	if err != nil {
		return err
	}

	fromUser.Experience -= xpAmount
	toUser.Experience += xpAmount

	if err := ctrl.UpdateUser(fromUser); err != nil {
		return err
	}
	if err := ctrl.UpdateUser(toUser); err != nil {
		return err
	}

	return nil
}

// SetMoney définit un montant d'argent pour un utilisateur
func (ctrl *UserController) SetMoney(userID string, moneyAmount int) error {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	user.Money = moneyAmount
	return ctrl.UpdateUser(user)
}

// SetAffinity définit un montant d'affinité pour un utilisateur
func (ctrl *UserController) SetAffinity(userID string, affinityAmount int) error {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	user.Affinity = affinityAmount
	return ctrl.UpdateUser(user)
}

// SetExperience définit un montant d'expérience pour un utilisateur
func (ctrl *UserController) SetExperience(userID string, xpAmount int) error {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	user.Experience = xpAmount
	return ctrl.UpdateUser(user)
}

// UpdateMoney met à jour le montant d'argent d'un utilisateur
func (ctrl *UserController) UpdateMoney(userID string, moneyAmount int) error {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	user.Money += moneyAmount
	return ctrl.UpdateUser(user)
}

// UpdateExperience met à jour le montant d'expérience d'un utilisateur
func (ctrl *UserController) UpdateExperience(userID string, xpAmount int) error {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return err
	}

	user.Experience += xpAmount
	return ctrl.UpdateUser(user)
}

// AddUserIfNotExists ajoute un utilisateur s'il n'existe pas déjà
func (ctrl *UserController) AddUserIfNotExists(userID, username string) error {
	exists, err := ctrl.UserExistsByDiscordID(userID)
	if err != nil {
		return err
	}
	if !exists {
		_, err := ctrl.CreateUser(userID, username, 0, 0, 0, "", 0, 0, 0, 0)
		return err
	}
	return nil
}

// GetMoney récupère le montant d'argent d'un utilisateur
func (ctrl *UserController) GetMoney(userID string) (int, error) {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return 0, err
	}
	return user.Money, nil
}

// GetExperience récupère le montant d'expérience d'un utilisateur
func (ctrl *UserController) GetExperience(userID string) (int, error) {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return 0, err
	}
	return user.Experience, nil
}

// GetAffinity récupère l'affinité d'un utilisateur
func (ctrl *UserController) GetAffinity(userID string) (int, error) {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return 0, err
	}
	return user.Affinity, nil
}

// GetScore récupère le score global d'un utilisateur (argent, expérience, affinité)
func (ctrl *UserController) GetScore(userID string) (int, error) {
	user, err := ctrl.GetUserByDiscordID(userID)
	if err != nil {
		return 0, err
	}
	// Exemple simple de calcul de score global (peut être modifié selon ta logique)
	totalScore := user.Money + user.Experience + user.Affinity
	return totalScore, nil
}
