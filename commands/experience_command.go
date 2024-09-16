package commands

import (
	"fmt"

	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// ExperienceCommand affiche l'affinité actuelle d'un utilisateur
func ExperienceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Empêcher le bot de répondre à lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	exp, exists := features.GetExperience(m.Author.ID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas encore gagné d'expérience.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez %d points d'expérience !", exp))
}
