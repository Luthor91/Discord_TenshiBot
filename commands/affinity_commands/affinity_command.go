package affinity_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// AffinityCommand gère les opérations d'affinité pour les utilisateurs
func AffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := fmt.Sprintf("%saffinity", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?affinity [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g]")
		return
	}

	targetUserID, affinityAmount, action, err := parseAffinityArgs(args, m, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	switch action {
	case "remove":
		handleRemoveAffinity(s, m, targetUserID, affinityAmount)
	case "set":
		handleSetAffinity(s, m, targetUserID, affinityAmount)
	case "add":
		handleAddAffinity(s, m, targetUserID, affinityAmount)
	case "get":
		handleGetAffinity(s, m, targetUserID)
	default:
		s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, ou -g.")
	}
}
