package services

import (
	"fmt"
	"log"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
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

	// Vérifie si l'utilisateur existe déjà
	if err := database.DB.Where("user_discord_id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Si l'utilisateur n'existe pas, vous pouvez décider de créer un nouvel utilisateur ou retourner une erreur
			return fmt.Errorf("utilisateur non trouvé avec ID Discord: %s", userID)
		}
		return err // Autre erreur
	}

	// Ajoute l'expérience
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
