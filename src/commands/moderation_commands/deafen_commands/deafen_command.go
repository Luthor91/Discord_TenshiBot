package deafen_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// DeafenCommand gère la commande ?deafen
// Usage : ?deafen <user> <durée> [raison]
func DeafenCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "deafen"

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// ?deafen
	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 3 {
		showHelpMessage(s, m.ChannelID)
		return
	}

	target := discord.HandleTarget(s, m, parts[1])
	if target == nil {
		return
	}

	duration, err := discord.ParseDuration(parts[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Durée invalide (ex: 10m, 1h)")
		return
	}

	// Raison optionnelle
	reason := "Aucune raison spécifiée"
	if len(parts) > 3 {
		reason = strings.Join(parts[3:], " ")
	}

	discord.DeafenUser(s, m, target.ID, duration, reason)

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			"Utilisateur %s rendu sourd pour %v : %s",
			target.Username,
			duration,
			reason,
		),
	)
}
