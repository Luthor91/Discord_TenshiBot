package commands

import (
	"fmt"

	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// AffinityCommand affiche l'affinité actuelle d'un utilisateur
func AffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Empêcher le bot de répondre à lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Récupérer l'affinité de l'utilisateur
	user, exists := features.GetUserAffinity(m.Author.ID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Aucune affinité trouvée pour cet utilisateur.")
		return
	}

	// Envoyer l'affinité de l'utilisateur
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité actuelle de %s : %d", user.Username, user.Affinity))
}
