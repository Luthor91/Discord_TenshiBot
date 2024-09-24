package services

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

// UserService est un service pour gérer les opérations liées aux utilisateurs
type UserService struct {
	userCtrl *controllers.UserController
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService() *UserService {
	return &UserService{
		userCtrl: controllers.NewUserController(),
	}
}

////////// GESTION DES UTILISATEURS //////////

// AddUserIfNotExists ajoute un utilisateur à la base de données s'il n'existe pas déjà
func (service *UserService) AddUserIfNotExists(userDiscordID, username string) error {
	return service.userCtrl.AddUserIfNotExists(userDiscordID, username)
}

// GetAllUsers utilise le UserController pour récupérer tous les utilisateurs
func (service *UserService) GetAllUsers() ([]models.User, error) {
	return service.userCtrl.GetAllUsers()
}

// GetUserByDiscordID utilise le UserController pour récupérer un utilisateur par son ID Discord
func (service *UserService) GetUserByDiscordID(userDiscordID string) (*models.User, error) {
	return service.userCtrl.GetUserByDiscordID(userDiscordID)
}

// GetUserRankAndScoreByCategory renvoie le rang et le score d'un utilisateur pour une catégorie spécifique
func (service *UserService) GetUserRankAndScoreByCategory(userDiscordID string, category string) (int, int, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return 0, 0, err
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
			return (u.Money + u.Experience + u.Affinity) / 3
		}
	default:
		return 0, 0, errors.New("invalid category")
	}

	for rank, user := range users {
		if user.UserDiscordID == userDiscordID {
			return rank + 1, scoreFunc(user), nil
		}
	}
	return 0, 0, errors.New("user not found")
}

// GetUserRankAndScore renvoie le rang et le score combiné d'un utilisateur
func (service *UserService) GetUserRankAndScore(userDiscordID string) (int, int, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return 0, 0, err
	}
	for rank, user := range users {
		if user.UserDiscordID == userDiscordID {
			return rank + 1, user.Money + user.Experience + user.Affinity, nil
		}
	}
	return 0, 0, errors.New("user not found")
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
func (service *UserService) CanReceiveDailyReward(user *models.User) (bool, time.Duration, error) {
	if user.LastDailyReward == "" {
		return true, 0, nil
	}

	lastRewardTime, err := time.Parse(time.RFC3339, user.LastDailyReward)
	if err != nil {
		return true, 0, err
	}

	now := time.Now()
	if now.Sub(lastRewardTime).Hours() >= 24 {
		return true, 0, nil
	}
	return false, time.Until(lastRewardTime.Add(24 * time.Hour)), nil
}

// UpdateDailyMoney met à jour la monnaie quotidienne et la date de la dernière récompense
func (service *UserService) UpdateDailyMoney(user *models.User, amount int) error {
	user.Money += amount
	user.LastDailyReward = time.Now().Format(time.RFC3339)
	return service.userCtrl.UpdateUser(user)
}

////////// GESTION DE L'EXPERIENCE ET AFFINITÉ //////////

// AddExperience ajoute de l'expérience à un utilisateur
func (service *UserService) AddExperience(user *models.User, amount int) error {
	user.Experience += amount
	return service.userCtrl.UpdateUser(user)
}

// SetExperience définit la quantité d'expérience pour un utilisateur
func (service *UserService) SetExperience(userDiscordID string, xp int) error {
	return service.userCtrl.SetExperience(userDiscordID, xp)
}

// GetExperience renvoie l'expérience d'un utilisateur
func (service *UserService) GetExperience(userDiscordID string) (int, error) {
	return service.userCtrl.GetExperience(userDiscordID)
}

// SetAffinity met à jour l'affinité d'un utilisateur
func (service *UserService) SetAffinity(userDiscordID string, affinityAmount int) error {
	return service.userCtrl.SetAffinity(userDiscordID, affinityAmount)
}

////////// GESTION DE LA MONNAIE //////////

// AddMoney ajoute de la monnaie à un utilisateur
func (service *UserService) AddMoney(user *models.User, amount int) error {
	user.Money += amount
	return service.userCtrl.UpdateUser(user)
}

// GetMoney renvoie la quantité de monnaie d'un utilisateur
func (service *UserService) GetMoney(userDiscordID string) (int, error) {
	return service.userCtrl.GetMoney(userDiscordID)
}

// UpdateMoney ajoute de l'argent à l'utilisateur
func (service *UserService) UpdateMoney(userDiscordID string, moneyAmount int) error {
	return service.userCtrl.UpdateMoney(userDiscordID, moneyAmount)
}

// SetMoney met à jour la quantité d'argent d'un utilisateur
func (service *UserService) SetMoney(userDiscordID string, moneyAmount int) error {
	return service.userCtrl.SetMoney(userDiscordID, moneyAmount)
}

////////// TRANSFERTS ENTRE UTILISATEURS //////////

// GiveExperience transfère une quantité d'Experience d'un utilisateur à un autre
func (service *UserService) GiveExperience(fromUserDiscordID, toUserDiscordID string, xpAmount int) error {
	return service.userCtrl.GiveExperience(fromUserDiscordID, toUserDiscordID, xpAmount)
}

// GiveMoney transfère une quantité d'argent d'un utilisateur à un autre
func (service *UserService) GiveMoney(fromUserDiscordID, toUserDiscordID string, moneyAmount int) error {
	return service.userCtrl.GiveMoney(fromUserDiscordID, toUserDiscordID, moneyAmount)
}

////////// CALCUL DU SCORE GENERAL //////////

// GetScore renvoie le score total pour un utilisateur
func (service *UserService) GetScore(userDiscordID string) (int, error) {
	return service.userCtrl.GetScore(userDiscordID)
}

// GetAffinity renvoie l'affinité pour un utilisateur
func (service *UserService) GetAffinity(userDiscordID string) (int, error) {
	return service.userCtrl.GetAffinity(userDiscordID)
}

// UserApplyEffects vérifie l'effet à appliquer sur l'utilisateur, ici uniquement le timeout
func (service *UserService) UserApplyEffects(s *discordgo.Session, guildID string, target *models.User) error {
	// Vérifier si l'utilisateur a un timeout à appliquer
	if target.TimeoutEnd.IsZero() {
		// Pas de timeout à appliquer
		return nil
	}

	// Calculer la durée restante avant la fin du timeout
	duration := time.Until(target.TimeoutEnd)

	// Si la durée est positive (timeout encore valide)
	if duration > 0 {
		// Appliquer le timeout avec la durée restante
		err := discord.TimeoutUser(s, guildID, target.UserDiscordID, duration)
		if err != nil {
			return fmt.Errorf("erreur lors de l'application du timeout : %v", err)
		}
		log.Printf("Timeout appliqué à l'utilisateur %s pour une durée de %s", target.UserDiscordID, duration)
	} else {
		log.Printf("Le timeout de l'utilisateur %s est expiré, aucune action n'a été prise.", target.UserDiscordID)
	}

	return nil
}
