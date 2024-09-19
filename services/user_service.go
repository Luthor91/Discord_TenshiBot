package services

import (
	"errors"
	"sort"
	"sync"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// UserService est un service pour gérer les opérations liées aux utilisateurs
type UserService struct {
	db *gorm.DB
	mu sync.Mutex // Ajout du mutex pour la gestion des accès concurrents
}

// NewUserService crée une nouvelle instance de UserService
func NewUserService() *UserService {
	return &UserService{
		db: database.DB,
	}
}

func (service *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := service.db.Order("money + experience + affinity DESC").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserRankAndScoreByCategory renvoie le rang et le score d'un utilisateur pour une catégorie spécifique
func (service *UserService) GetUserRankAndScoreByCategory(userID string, category string) (int, int, bool, error) {
	users, err := service.GetAllUsers()
	if err != nil {
		return 0, 0, false, err
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
	var users []models.User
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

	if err := service.db.Find(&users).Error; err != nil {
		return nil, err
	}

	sort.Slice(users, func(i, j int) bool {
		return scoreFunc(users[i]) > scoreFunc(users[j])
	})
	return users, nil
}

// GetUserScore calcule le score d'un utilisateur en fonction de la catégorie spécifiée
func (service *UserService) GetUserScore(user models.User, category string) int {
	switch category {
	case "money":
		return user.Money
	case "affinity":
		return user.Affinity
	case "xp":
		return user.Experience
	case "general":
		return (user.Money + user.Affinity + user.Experience) / 3
	default:
		return 0
	}
}

// AddUserIfNotExists ajoute un utilisateur à la base de données s'il n'existe pas déjà
func (service *UserService) AddUserIfNotExists(userID, username string) error {
	var user models.User
	if err := service.db.Where("user_discord_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := models.User{
				UserDiscordID: userID,
				Username:      username,
				Money:         0,
				Experience:    0,
				Affinity:      0,
			}
			if err := service.db.Create(&newUser).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}
