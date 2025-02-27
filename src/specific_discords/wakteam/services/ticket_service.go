package services

import (
	"fmt"

	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/controllers"
	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/models"
)

// TicketService contient la logique pour gérer les tickets.
type TicketService struct {
	Controller *controllers.TicketController
}

// NewTicketService crée une nouvelle instance de TicketService.
func NewTicketService() *TicketService {
	return &TicketService{Controller: controllers.NewTicketController()}
}

// CreateTicket crée un nouveau ticket pour un utilisateur.
func (s *TicketService) CreateTicket(author, content string) error {
	ticket := &models.Ticket{
		Author:  author,
		Content: content,
		Status:  "en_cours",
	}

	if err := s.Controller.CreateTicket(ticket); err != nil {
		return fmt.Errorf("échec de la création du ticket : %w", err)
	}
	return nil
}

// ResolveTicket met le dernier ticket d'un utilisateur en "validé".
func (s *TicketService) ResolveTicket(author string) error {
	ticket, err := s.Controller.GetLastTicket(author)
	if err != nil {
		return fmt.Errorf("échec de la récupération du dernier ticket : %w", err)
	}

	if err := s.Controller.UpdateTicketStatus(ticket.ID, "valide"); err != nil {
		return fmt.Errorf("échec de la validation du ticket : %w", err)
	}
	return nil
}

// RejectTicket met le dernier ticket d'un utilisateur en "refusé".
func (s *TicketService) RejectTicket(author string) error {
	ticket, err := s.Controller.GetLastTicket(author)
	if err != nil {
		return fmt.Errorf("échec de la récupération du dernier ticket : %w", err)
	}

	if err := s.Controller.UpdateTicketStatus(ticket.ID, "refuse"); err != nil {
		return fmt.Errorf("échec du refus du ticket : %w", err)
	}
	return nil
}

// DeleteTicket supprime un ticket en fonction de son ID.
func (s *TicketService) DeleteTicket(ticketID uint) error {
	if err := s.Controller.DeleteTicket(ticketID); err != nil {
		return fmt.Errorf("échec de la suppression du ticket : %w", err)
	}
	return nil
}
