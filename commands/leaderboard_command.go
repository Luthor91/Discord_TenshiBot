package commands

import (
	"fmt"
	"sort"

	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// Leaderboard affiche le classement des utilisateurs en fonction de leur monnaie
func LeaderboardCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Récupérer la monnaie de tous les utilisateurs
	users := features.GetAllUsersMoney()

	// Trier les utilisateurs par monnaie décroissante
	sort.Slice(users, func(i, j int) bool {
		return users[i].Money > users[j].Money
	})

	// Construire le message de réponse
	response := "Classement des utilisateurs :\n"
	for i, user := range users {
		// Récupérer le nom d'utilisateur via l'API Discord
		member, err := s.GuildMember(m.GuildID, user.UserID)
		if err != nil {
			response += fmt.Sprintf("%d. Utilisateur %s - %d pièces\n", i+1, user.UserID, user.Money)
			continue
		}
		response += fmt.Sprintf("%d. %s - %d pièces\n", i+1, member.User.Username, user.Money)
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
