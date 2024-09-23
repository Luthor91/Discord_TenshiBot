package services

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Luthor91/Tenshi/controllers" // Importer le contrôleur
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// UserService est un service pour gérer les opérations liées aux utilisateurs
type UserService struct {
	userCtrl *controllers.UserController
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService(controller *controllers.UserController) *UserService {
	return &UserService{
		userCtrl: controller,
	}
}

////////// GESTION DES UTILISATEURS //////////

// AddUserIfNotExists ajoute un utilisateur à la base de données s'il n'existe pas déjà
func (service *UserService) AddUserIfNotExists(userID, username string) error {
	return service.userCtrl.AddUserIfNotExists(userID, username)
}

// GetAllUsers utilise le UserController pour récupérer tous les utilisateurs
func (service *UserService) GetAllUsers() ([]models.User, error) {
	return service.userCtrl.GetAllUsers()
}

// GetUserByID utilise le UserController pour récupérer un utilisateur par ID
func (service *UserService) GetUserByID(userID uint) (*models.User, error) {
	userPtr, err := service.userCtrl.GetUserByID(userID)
	if err != nil {
		return &models.User{}, err // Gérer l'erreur si nécessaire
	}
	return userPtr, nil // Déférencer le pointeur pour retourner un modèle User
}

// GetUserByDiscordID utilise le UserController pour récupérer un utilisateur par ID Discord
func (service *UserService) GetUserByDiscordID(userDiscordID string) (*models.User, error) {
	userPtr, err := service.userCtrl.GetUserByDiscordID(userDiscordID)
	if err != nil {
		return &models.User{}, err // Gérer l'erreur si nécessaire
	}
	return userPtr, nil // Déférencer le pointeur pour retourner un modèle User
}

// GetUserRankAndScoreByCategory renvoie le rang et le score d'un utilisateur pour une catégorie spécifique
func (service *UserService) GetUserRankAndScoreByCategory(userID string, category string) (int, int, bool, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return 0, 0, false, err
	}

	// Fonction de score selon la catégorie
	var scoreFunc func(models.User) int
	switch category {
	case "money":
		scoreFunc = func(u models.User) int { return u.Money }
	case "affinity":
		scoreFunc = func(u models.User) int { return u.Affinity }
	case "xp":
		scoreFunc = func(u models.User) int { return u.Experience }
	case "general":
		scoreFunc = func(u models.User) int {
			return (u.Money + u.Experience + u.Affinity) / 3
		}
	default:
		return 0, 0, false, errors.New("invalid category")
	}

	for rank, user := range users {
		if user.UserDiscordID == userID {
			return rank + 1, scoreFunc(user), true, nil
		}
	}
	return 0, 0, false, nil
}

// GetUserRankAndScore renvoie le rang et le score combiné d'un utilisateur
func (service *UserService) GetUserRankAndScore(userID string) (int, int, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return 0, 0, err
	}
	for rank, user := range users {
		if user.UserDiscordID == userID {
			return rank + 1, user.Money + user.Experience + user.Affinity, nil
		}
	}
	return 0, 0, nil
}

// GetAllUsersByCategory renvoie tous les utilisateurs triés par une catégorie spécifique
func (service *UserService) GetAllUsersByCategory(category string) ([]models.User, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var scoreFunc func(models.User) int
	switch category {
	case "money":
		scoreFunc = func(u models.User) int { return u.Money }
	case "affinity":
		scoreFunc = func(u models.User) int { return u.Affinity }
	case "xp":
		scoreFunc = func(u models.User) int { return u.Experience }
	case "general":
		scoreFunc = func(u models.User) int {
			return u.Money + u.Experience + u.Affinity
		}
	default:
		return nil, errors.New("invalid category")
	}

	sort.Slice(users, func(i, j int) bool {
		return scoreFunc(users[i]) > scoreFunc(users[j])
	})
	return users, nil
}

////////// GESTION DES RÉCOMPENSES QUOTIDIENNES //////////

// CanReceiveDailyReward vérifie si l'utilisateur peut recevoir une récompense quotidienne
func (service *UserService) CanReceiveDailyReward(user *models.User) (bool, time.Duration) {
	if user.LastDailyReward == "" {
		return true, 0
	}

	lastRewardTime, err := time.Parse(time.RFC3339, user.LastDailyReward)
	if err != nil {
		return true, 0
	}

	now := time.Now()
	if now.Sub(lastRewardTime).Hours() >= 24 {
		return true, 0
	}
	return false, time.Until(lastRewardTime.Add(24 * time.Hour))
}

// UpdateDailyMoney met à jour la monnaie quotidienne et la date de la dernière récompense
func (service *UserService) UpdateDailyMoney(user *models.User, amount int) error {
	user.Money += amount
	user.LastDailyReward = time.Now().Format(time.RFC3339)
	return service.userCtrl.UpdateUser(user)
}

////////// GESTION DE L'EXPÉRIENCE ET AFFINITÉ //////////

// AddExperience ajoute de l'expérience à un utilisateur
func (service *UserService) AddExperience(user *models.User, amount int) error {
	user.Experience += amount
	return service.userCtrl.UpdateUser(user)
}

// GetExperience renvoie l'expérience d'un utilisateur
func (service *UserService) GetExperience(userID string) (int, bool) {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur %s: %v", userID, err)
		return 0, false
	}
	return user.Experience, true
}

// SetExperience définit la quantité d'expérience pour un utilisateur
func (service *UserService) SetExperience(userID string, xp int) error {
	var user models.User

	if err := database.DB.Where("user_discord_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("utilisateur non trouvé avec ID Discord: %s", userID)
		}
		return err
	}

	user.Experience = xp
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	log.Printf("Expérience mise à jour pour l'utilisateur %s : %d", userID, xp)
	return nil
}

// SetAffinity met à jour la quantité d'affinité d'un utilisateur
func (service *UserService) SetAffinity(userID string, affinityAmount int) error {
	return service.userCtrl.SetAffinity(userID, affinityAmount)
}

////////// GESTION DE LA MONNAIE //////////

// AddMoney ajoute de la monnaie à un utilisateur
func (service *UserService) AddMoney(user *models.User, amount int) error {
	user.Money += amount
	return service.userCtrl.UpdateUser(user)
}

// GetMoney renvoie la quantité de monnaie d'un utilisateur
func (service *UserService) GetMoney(discordID string) (int, error) {
	user, err := service.userCtrl.GetUserByDiscordID(discordID)
	if err != nil {
		return 0, err
	}
	return user.Money, nil
}

// UpdateMoney ajoute de l'argent à l'utilisateur
func (service *UserService) UpdateMoney(userID string, moneyAmount int) error {
	return service.userCtrl.UpdateMoney(userID, moneyAmount)
}

// SetMoney met à jour la quantité d'argent d'un utilisateur
func (service *UserService) SetMoney(userID string, moneyAmount int) error {
	return service.userCtrl.SetMoney(userID, moneyAmount)
}

////////// TRANSFERTS ENTRE UTILISATEURS //////////

// GiveXP transfère une quantité d'XP d'un utilisateur à un autre
func (service *UserService) GiveXP(fromUserID, toUserID string, xpAmount int) error {
	return service.userCtrl.GiveXP(fromUserID, toUserID, xpAmount)
}

// GiveMoney transfère une quantité d'argent d'un utilisateur à un autre
func (service *UserService) GiveMoney(fromUserID, toUserID string, moneyAmount int) error {
	return service.userCtrl.GiveMoney(fromUserID, toUserID, moneyAmount)
}

// //////// CALCUL DU SCORE GENERAL //////////

// GetUserScore renvoie le score total (somme de l'argent, de l'affinité et de l'expérience) d'un utilisateur
func (service *UserService) GetScore(userID string) (int, error) {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la récupération de l'utilisateur : %v", err)
	}

	// Calculer le score total en additionnant l'argent, l'affinité et l'expérience
	totalScore := user.Money + user.Affinity + user.Experience
	return totalScore, nil
}
