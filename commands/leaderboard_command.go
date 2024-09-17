package commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

// LeaderboardCommand affiche le classement des utilisateurs en fonction de la monnaie, affinité et expérience
func LeaderboardCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sleaderboard", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		args := strings.Split(m.Content, " ")
		category := "general"
		if len(args) >= 2 {
			category = args[1]
		}

		var users []models.User

		// Utiliser les fonctions spécifiques pour chaque catégorie
		switch category {
		case "money":
			users = features.GetUsersByMoney()
		case "affinity":
			users = features.GetUsersByAffinity()
		case "xp":
			users = features.GetUsersByXP()
		case "general":
			users = features.GetUsersByGeneral()
		default:
			s.ChannelMessageSend(m.ChannelID, "Type de classement invalide. Choisissez parmi money, affinity, xp, ou general.")
			return
		}

		// Construire le message de réponse
		response := fmt.Sprintf("Classement %s des utilisateurs :\n", category)
		for i, user := range users {
			member, err := s.GuildMember(m.GuildID, user.UserID)
			if err != nil {
				response += fmt.Sprintf("%d. Utilisateur %s - %d\n", i+1, user.UserID, features.GetUserScore(user, category))
				continue
			}
			response += fmt.Sprintf("%d. %s - %d\n", i+1, member.User.Username, features.GetUserScore(user, category))
		}

		s.ChannelMessageSend(m.ChannelID, response)
	}
}
