package controllers

import (
	"fmt"

	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/models"
	"gorm.io/gorm"
)

// TicketController gère les opérations sur les tickets.
type TicketController struct {
	DB *gorm.DB
}

// NewTicketController crée une instance de TicketController.
func NewTicketController() *TicketController {
	return &TicketController{DB: database.DB}
}

// CreateTicket ajoute un ticket à la base de données.
func (c *TicketController) CreateTicket(ticket *models.Ticket) error {
	if err := c.DB.Create(ticket).Error; err != nil {
		return fmt.Errorf("échec de la création du ticket : %w", err)
	}
	return nil
}

// UpdateTicketStatus met à jour le statut d'un ticket.
func (c *TicketController) UpdateTicketStatus(ticketID uint, status string) error {
	if err := c.DB.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("status", status).Error; err != nil {
		return fmt.Errorf("échec de la mise à jour du statut du ticket : %w", err)
	}
	return nil
}

// GetLastTicket récupère le dernier ticket d'un utilisateur.
func (c *TicketController) GetLastTicket(author string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := c.DB.Where("author = ?", author).Order("created_at desc").First(&ticket).Error; err != nil {
		return nil, fmt.Errorf("aucun ticket trouvé : %w", err)
	}
	return &ticket, nil
}

// DeleteTicket supprime un ticket en fonction de son ID.
func (c *TicketController) DeleteTicket(ticketID uint) error {
	var ticket models.Ticket
	if err := c.DB.First(&ticket, ticketID).Error; err != nil {
		return fmt.Errorf("ticket introuvable : %w", err)
	}
	return c.DB.Delete(&ticket).Error
}
