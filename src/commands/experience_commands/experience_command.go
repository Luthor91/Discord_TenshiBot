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

	command := fmt.Sprintf("%sexperience", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%sxp", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	args := strings.Fields(m.Content)

	// Affiche l'aide si -h est spécifié
	if len(args) > 1 && args[1] == "-h" {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m]")
		return
	}

	// Si aucun argument n'est fourni, affiche l'XP de l'utilisateur qui a exécuté la commande
	if len(args) < 2 {
		targetUserID := m.Author.ID // Définit l'utilisateur par défaut
		handleGetXP(s, m, targetUserID)
		return
	}

	// Parse les arguments pour obtenir l'utilisateur cible, la quantité d'XP, et l'action
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
