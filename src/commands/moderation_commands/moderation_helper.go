package moderation_commands

import (
	"errors"
	"fmt"

	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// ParseArgs parses the command line arguments
func parseArgs(args []string) (map[string]string, error) {
	parsedArgs := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n", "-b", "-w", "-k", "-m", "-d", "-to", "-mv":
			if i+1 < len(args) {
				parsedArgs[arg] = args[i+1]
				i++
			} else {
				return nil, errors.New("missing value for " + arg)
			}
		case "-t":
			if i+1 < len(args) {
				parsedArgs[arg] = args[i+1]
				i++
			} else {
				return nil, errors.New("missing duration value for -t")
			}
		case "-r", "-rw":
			parsedArgs[arg] = "true" // Just a flag, doesn't need a value
		}
	}
	return parsedArgs, nil
}

func warnUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	// Appeler la méthode AddWarn du WarnService pour ajouter un avertissement
	warnService := services.NewWarnService(s, m.GuildID)
	err := warnService.AddWarn(userID, reason, m.Author.ID) // m.Author.ID pour l'admin qui donne l'avertissement
	if err != nil {
		// Gérer l'erreur et envoyer un message dans le canal
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de l'envoi de l'avertissement à l'utilisateur %s : %v", userID, err))
		return
	}

	// Si tout se passe bien, envoyer un message de confirmation
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissement envoyé à l'utilisateur %s : %s", userID, reason))
}

func resetAllUserStatus(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	// Implement reset logic to clear all punishments (ban, mute, deafen, etc.)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Statut de l'utilisateur %s réinitialisé.", userID))
}

func resetUserWarnings(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	// Réinitialiser les avertissements de l'utilisateur
	warnService := services.NewWarnService(s, m.GuildID)

	err := warnService.ResetWarns(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la réinitialisation des avertissements de l'utilisateur %s: %v", userID, err))
		return
	}

	// Confirmer que les avertissements ont été réinitialisés
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissements de l'utilisateur %s réinitialisés.", userID))
}
