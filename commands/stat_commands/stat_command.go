package stat_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config" // Assurez-vous de pointer vers votre fichier de configuration
	"github.com/bwmarrin/discordgo"
)

// StatCommand gère les différentes options pour afficher les statistiques du bot, du serveur et de l'utilisateur.
func StatCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Vérifie que le message ne provient pas du bot lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifie que le message commence par le préfixe de commande
	command := fmt.Sprintf("%sstat", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Parsing des arguments
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier une option: -u (utilisateur), -s (serveur), ou -b (bot).")
		return
	}

	// Gère les différentes options
	option := args[1]
	switch option {
	case "-u":
		userStatsCommand(s, m)
	case "-s":
		serverStatsCommand(s, m)
	case "-b":
		botStatsCommand(s, m)
	case "-p":
		botPerfsCommand(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Option inconnue. Utilisez -u (utilisateur), -s (serveur), ou -b (bot).")
	}
}
