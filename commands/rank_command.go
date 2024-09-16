package commands

import (
	"fmt"

	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// Rank affiche le rang d'un utilisateur dans le classement
func RankCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	rank, money, found := features.GetUserRankAndMoney(userID)
	if !found {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de votre rang.")
		return
	}

	message := fmt.Sprintf("%s, vous êtes classé %dème avec %d pièces.", m.Author.Username, rank, money)
	s.ChannelMessageSend(m.ChannelID, message)
}
