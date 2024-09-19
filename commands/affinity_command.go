package commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// AffinityCommand affiche l'affinité actuelle d'un utilisateur
func AffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%saffinity", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {

		// Récupérer l'affinité de l'utilisateur
		user, exists := services.GetUserAffinity(m.Author.ID)
		if !exists {
			return
		}

		// Envoyer l'affinité de l'utilisateur
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité actuelle de %s : %d", user.Username, user.Affinity))
	}
}
