package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Luthor91/Tenshi/controllers" // Remplacez par le chemin correct de votre package controllers
	"github.com/Luthor91/Tenshi/models"
)

// InvestmentService contient la logique pour gérer les investissements.
type InvestmentService struct {
	Controller *controllers.InvestmentController
}

// NewInvestmentService crée un nouveau service d'investissement.
func NewInvestmentService() *InvestmentService {

	return &InvestmentService{Controller: controllers.NewInvestmentController()}
}

// CreateInvestment crée un nouvel investissement pour un utilisateur.
func (s *InvestmentService) CreateInvestment(discordID string, amount int) error {
	user, err := NewUserService().GetUserByDiscordID(discordID)
	if err != nil {
		return fmt.Errorf("échec de la récupération de l'utilisateur : %w", err)
	}

	// Vérifier si l'utilisateur a suffisamment d'argent pour investir
	if user.Money < amount {
		return fmt.Errorf("vous n'avez pas assez d'argent pour effectuer cet investissement")
	}

	// Enregistrer l'investissement dans la base de données
	investment := models.Investment{
		UserID:    user.ID,
		Amount:    amount,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	err = s.Controller.CreateInvestment(&investment)
	if err != nil {
		return fmt.Errorf("erreur lors de l'enregistrement de l'investissement : %w", err)
	}

	// Déduire l'argent de l'utilisateur
	newMoneyAmount := user.Money - amount
	err = NewUserService().AddMoney(user.UserDiscordID, newMoneyAmount)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du solde : %w", err)
	}

	return nil
}

// CollectInvestment permet de récupérer l'investissement d'un utilisateur.
func (s *InvestmentService) CollectInvestment(discordID string) (int, float64, error) {
	user, err := NewUserService().GetUserByDiscordID(discordID)
	if err != nil {
		return 0, 0, fmt.Errorf("échec de la récupération de l'utilisateur : %w", err)
	}

	// Récupérer le dernier investissement de l'utilisateur
	lastInvestment, err := s.Controller.GetLastInvestmentUser(user.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("aucun investissement trouvé : %w", err)
	}

	// Vérifier si 24 heures se sont écoulées
	if time.Since(lastInvestment.CreatedAt).Hours() < 24 {
		return 0, 0, fmt.Errorf("vous devez attendre 24 heures avant de récupérer votre investissement")
	}

	// Générer un multiplicateur aléatoire (par exemple, entre 0.5 et 1.5)
	multiplier := 0.5 + rand.Float64()*(1.5-0.5)
	returnAmount := int(float64(lastInvestment.Amount) * multiplier)

	// Mettre à jour le solde de l'utilisateur
	newMoneyAmount := user.Money + returnAmount
	err = NewUserService().AddMoney(user.UserDiscordID, newMoneyAmount)
	if err != nil {
		return 0, 0, fmt.Errorf("erreur lors de la mise à jour du solde : %w", err)
	}

	// Supprimer l'investissement après récupération
	err = s.Controller.DeleteInvestment(lastInvestment.ID)
	if err != nil {
		return 0, 0, fmt.Errorf("erreur lors de la suppression de l'investissement : %w", err)
	}

	return returnAmount, multiplier, nil
}

// GetLastInvestmentUser récupère le dernier investissement d'un utilisateur par son ID Discord.
func (s *InvestmentService) GetLastInvestmentUser(discordID string) (*models.Investment, error) {
	// Appelle la méthode du contrôleur pour récupérer le dernier investissement
	user, _ := NewUserService().GetUserByDiscordID(discordID)
	investment, err := s.Controller.GetLastInvestmentUser(user.ID)
	if err != nil {
		return nil, err
	}
	return investment, nil
}

// DeleteInvestment supprime un investissement en fonction de son ID.
func (is *InvestmentService) DeleteInvestment(id uint) error {
	return is.Controller.DeleteInvestment(id)
}
