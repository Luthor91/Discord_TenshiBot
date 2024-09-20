package services

import (
	"fmt"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
)

// WarnService gère la logique de l'application pour les warns
type WarnService struct {
	warnController *controllers.WarnController
}

// NewWarnService crée un nouveau WarnService
func NewWarnService(warnController *controllers.WarnController) *WarnService {
	return &WarnService{warnController: warnController}
}

// AddWarn ajoute un nouveau warn à un utilisateur
func (ws *WarnService) AddWarn(userDiscordID, reason, adminID string) error {
	err := ws.warnController.CreateWarn(userDiscordID, reason, adminID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du warn: %v", err)
	}
	return nil
}

// GetWarns retourne la liste des warns d'un utilisateur
func (ws *WarnService) GetWarns(userDiscordID string) ([]models.Warn, error) {
	warns, err := ws.warnController.GetWarnsByUserDiscordID(userDiscordID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des warns: %v", err)
	}
	return warns, nil
}

// RemoveWarn supprime un warn
func (ws *WarnService) RemoveWarn(warnID uint) error {
	err := ws.warnController.DeleteWarn(warnID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du warn: %v", err)
	}
	return nil
}
