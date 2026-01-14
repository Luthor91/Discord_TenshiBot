package services

import (
	"github.com/Luthor91/DiscordBot/controllers"
	"github.com/Luthor91/DiscordBot/models"
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
