package stat_commands

import (
	"fmt"
	"strings"

	// Assurez-vous de pointer vers votre fichier de configuration

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// StatCommand gère les différentes options pour afficher les statistiques du bot, du serveur et de l'utilisateur.
func StatCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Vérifie que le message ne provient pas du bot lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifie que l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires.")
		return
	}

	// Définir le préfixe de commande
	command := fmt.Sprintf("%sstat", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer et analyser les arguments de la commande
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Vérifier qu'il y a des arguments
	if len(parsedArgs) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier une option : `-u` (utilisateur), `-s` (serveur), `-b` (bot), ou `-c` (canal).")
		return
	}

	// Gère les différentes options
	option := parsedArgs[0].Arg
	switch option {
	case "-u":
		discord.PrintUserStats(s, m)
	case "-s":
		discord.PrintServerStats(s, m)
	case "-b":
		discord.PrintBotStats(s, m)
	case "-c":
		discord.PrintChannelStats(s, m)
	default:
		s.ChannelMessageSend(m.ChannelID, "Option inconnue. Utilisez `-u` (utilisateur), `-s` (serveur), `-b` (bot), ou `-c` (canal).")
	}
}
