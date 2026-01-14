package services

import (
	"fmt"
	"time"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/controllers"
	"github.com/Luthor91/DiscordBot/models"
	"github.com/bwmarrin/discordgo"
)

// WarnService gère la logique de l'application pour les warns
type WarnService struct {
	warnController *controllers.WarnController
	discordSession *discordgo.Session // Ajoute la session Discord au service
	guildID        string             // Ajoute la guildID (ID du serveur)
}

// NewWarnService crée un nouveau WarnService avec la session Discord et l'ID du serveur
func NewWarnService(discordSession *discordgo.Session, guildID string) *WarnService {
	return &WarnService{
		warnController: controllers.NewWarnController(),
		discordSession: discordSession, // Ajoute la session Discord
		guildID:        guildID,        // Ajoute l'ID du serveur Discord
	}
}

// AddWarn ajoute un nouveau warn à un utilisateur
func (ws *WarnService) AddWarn(userDiscordID, reason, adminID string) (int64, error) {
	// Création du warn
	if err := ws.warnController.CreateWarn(userDiscordID, reason, adminID); err != nil {
		return 0, fmt.Errorf("erreur lors de l'ajout du warn: %v", err)
	}

	// Comptage
	warnCount, err := ws.warnController.CountWarnsByUser(userDiscordID)
	if err != nil {
		return 0, fmt.Errorf("erreur lors de la récupération du nombre de warns: %v", err)
	}

	// Timeout sur palier
	if warnCount > 0 && warnCount%3 == 0 {
		tier := warnCount / 3
		timeoutDuration := time.Duration(5*tier) * time.Minute

		if err := discord.TimeoutUser(
			ws.discordSession,
			ws.guildID,
			userDiscordID,
			timeoutDuration,
		); err != nil {
			return warnCount, fmt.Errorf("erreur lors de l'application du timeout, la cible est administrateur.")
		}
	}

	return warnCount, nil
}



// GetWarns retourne la liste des warns d'un utilisateur
func (ws *WarnService) GetWarns(userDiscordID string) ([]models.Warn, error) {
	warns, err := ws.warnController.GetWarnsByUserDiscordID(userDiscordID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des warns: %v", err)
	}
	return warns, nil
}

// ResetWarns réinitialise les avertissements d'un utilisateur
func (ws *WarnService) ResetWarns(userDiscordID string) error {
	err := ws.warnController.ResetWarns(userDiscordID)
	if err != nil {
		return fmt.Errorf("erreur lors de la réinitialisation des avertissements: %v", err)
	}
	return nil
}

// RemoveWarn supprime un warn
func (ws *WarnService) RemoveWarn(warnID uint) error {
	err := ws.warnController.DeleteWarn(warnID)
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du warn: %v", err)
	}
	return nil
}
