package features

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DeleteBanWordMessages supprime les messages contenant des mots interdits
func DeleteBanWordMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parcourir les mots interdits et vérifier s'ils sont présents dans le message
	for _, word := range badwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			// Supprimer le message
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Printf("Erreur lors de la suppression du message contenant un mot interdit: %v", err)
				return
			}
			break
		}
	}
}
