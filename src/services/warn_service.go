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
func (ws *WarnService) AddWarn(userDiscordID, reason, adminID string) error {
	// Ajouter un warn à l'utilisateur
	err := ws.warnController.CreateWarn(userDiscordID, reason, adminID)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du warn: %v", err)
	}

	// Récupérer le nombre total de warns pour cet utilisateur
	warnCount, err := ws.warnController.CountWarnsByUser(userDiscordID)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du nombre de warns: %v", err)
	}

	// Vérifier si le nombre de warns est un multiple de 3
	if warnCount%3 == 0 {
		// Calculer la durée du timeout en minutes (5 minutes par multiple de 3 warns)
		timeoutDuration := time.Duration(5*(warnCount/3)) * time.Minute

		// Appliquer un timeout à l'utilisateur sur Discord
		err := discord.TimeoutUser(ws.discordSession, ws.guildID, userDiscordID, timeoutDuration)
		if err != nil {
			return fmt.Errorf("erreur lors de l'application du timeout: %v", err)
		}

		// Retourner une confirmation que l'utilisateur a été mis en timeout
		fmt.Printf("L'utilisateur %s a été mis en timeout pour %v minutes en raison de %d warns.\n", userDiscordID, timeoutDuration.Minutes(), warnCount)
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
