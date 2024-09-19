package services

import (
	"log"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
)

// ExperienceService est un service pour gérer l'expérience des utilisateurs
type ExperienceService struct {
	userCtrl *controllers.UserController
}

// NewExperienceService crée une nouvelle instance de ExperienceService
func NewExperienceService() *ExperienceService {
	return &ExperienceService{
		userCtrl: &controllers.UserController{DB: database.DB},
	}
}

// AddExperience ajoute de l'expérience à un utilisateur et retourne une erreur si nécessaire
func (service *ExperienceService) AddExperience(userID string, amount int) error {
	var user models.User
	// Utilisation de .First() avec Where est correcte ici.
	if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return err
	}

	user.Experience += amount
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetExperience renvoie l'expérience d'un utilisateur
func (service *ExperienceService) GetExperience(userID string) (int, bool) {
	user, err := service.userCtrl.GetUserByDiscordID(userID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur %s: %v", userID, err)
		return 0, false
	}
	return user.Experience, true
}
