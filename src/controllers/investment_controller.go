package controllers

import (
	"fmt"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// InvestmentController gère les opérations sur les investissements dans la base de données.
type InvestmentController struct {
	DB *gorm.DB
}

// NewInvestmentController crée un nouveau contrôleur d'investissement.
func NewInvestmentController() *InvestmentController {
	return &InvestmentController{DB: database.DB}
}

// CreateInvestment crée un nouvel investissement dans la base de données.
func (c *InvestmentController) CreateInvestment(investment *models.Investment) error {
	if err := c.DB.Create(investment).Error; err != nil {
		return fmt.Errorf("échec de la création de l'investissement : %w", err)
	}
	return nil
}

// GetLastInvestmentUser récupère le dernier investissement d'un utilisateur.
func (c *InvestmentController) GetLastInvestmentUser(userID uint) (*models.Investment, error) {
	var investment models.Investment
	if err := c.DB.Where("user_id = ?", userID).Order("created_at desc").First(&investment).Error; err != nil {
		return nil, fmt.Errorf("aucun investissement trouvé : %w", err)
	}
	return &investment, nil
}

// DeleteInvestment supprime un investissement en fonction de son ID.
func (ic *InvestmentController) DeleteInvestment(id uint) error {
	var investment models.Investment
	if err := ic.DB.First(&investment, id).Error; err != nil {
		return err
	}
	return ic.DB.Delete(&investment).Error
}
