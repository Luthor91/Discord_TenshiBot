package ranking_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/models"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// LeaderboardCommand affiche le classement des utilisateurs en fonction de la monnaie, affinité et expérience
func LeaderboardCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	command := fmt.Sprintf("%sleaderboard", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%slb", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}
	args := strings.Split(m.Content, " ")
	category := "general"
	if len(args) >= 2 {
		category = args[1]
	}

	var users []models.User

	// Utiliser les fonctions spécifiques pour chaque catégorie
	switch category {
	case "money":
		users, _ = services.NewUserService().GetAllUsersByCategory("money")
	case "affinity":
		users, _ = services.NewUserService().GetAllUsersByCategory("affinity")
	case "xp":
		users, _ = services.NewUserService().GetAllUsersByCategory("xp")
	case "general":
		users, _ = services.NewUserService().GetAllUsersByCategory("general")
	default:
		s.ChannelMessageSend(m.ChannelID, "Type de classement invalide. Choisissez parmi money, affinity, xp, ou general.")
		return
	}

	// Construire le message de réponse
	response := fmt.Sprintf("Classement %s des utilisateurs :\n", category)
	for i, user := range users {
		member, err := s.GuildMember(m.GuildID, user.UserDiscordID)
		if err != nil {
			continue
		}
		score, err := services.NewUserService().GetScore(user.UserDiscordID)
		response += fmt.Sprintf("%d. %s - %d\n", i+1, member.User.Username, score)
	}

	s.ChannelMessageSend(m.ChannelID, response)
}
