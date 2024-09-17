package commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// ExperienceCommand affiche l'affinité actuelle d'un utilisateur
func ExperienceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sxp", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		exp, exists := features.GetExperience(m.Author.ID)
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas encore gagné d'expérience.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez %d points d'expérience !", exp))
	}
}
