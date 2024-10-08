package ranking_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// RankCommand affiche le rang d'un utilisateur dans un classement spécifique
func RankCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%srank", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		args := strings.Split(m.Content, " ")
		category := "general"
		if len(args) >= 2 {
			category = args[1]
		}

		userID := m.Author.ID

		var rank int
		var score int

		switch category {
		case "money":
			rank, score, _ = services.NewUserService().GetUserRankAndScoreByCategory(userID, "money")
		case "affinity":
			rank, score, _ = services.NewUserService().GetUserRankAndScoreByCategory(userID, "affinity")
		case "xp":
			rank, score, _ = services.NewUserService().GetUserRankAndScoreByCategory(userID, "xp")
		case "general":
			rank, score, _ = services.NewUserService().GetUserRankAndScoreByCategory(userID, "general")
		default:
			s.ChannelMessageSend(m.ChannelID, "Type de classement invalide. Choisissez parmi money, affinity, xp, ou general.")
			return
		}

		message := fmt.Sprintf("%s, vous êtes classé %dème avec %d.", m.Author.Username, rank, score)
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
