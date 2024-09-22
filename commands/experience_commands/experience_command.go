package experience_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// XPCommand gère les opérations d'expérience pour les utilisateurs
func ExperienceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := fmt.Sprintf("%sxp", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m]")
		return
	}

	targetUserID, xpAmount, action, err := parseXPArgs(args, m, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch action {
	case "remove":
		handleRemoveXP(s, m, targetUserID, xpAmount)
	case "set":
		handleSetXP(s, m, targetUserID, xpAmount)
	case "add":
		handleAddXP(s, m, targetUserID, xpAmount)
	case "me":
		handleGetXP(s, m, targetUserID)
	case "give":
		handleGiveXP(s, m, m.Author.ID, targetUserID, xpAmount) // Giver ID is the command author
	default:
		s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, -m, ou -h.")
	}
}
