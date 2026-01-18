package timeout_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// TimeoutCommand gère la commande ?timeout
// Usage : ?timeout <user> <durée> [raison]
func TimeoutCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "timeout"
	aliasCommand := config.AppConfig.BotPrefix + "to"

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, aliasCommand) {
		return
	}

	// ?timeout
	if m.Content == command || m.Content == aliasCommand {
		showHelpMessage(s, m.ChannelID)
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 3 {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Target
	target := discord.HandleTarget(s, m, parts[1])
	if target == nil {
		return
	}

	// Durée
	duration, err := discord.ParseDuration(parts[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Durée invalide (ex: 10m, 1h, 1d)")
		return
	}

	// Raison optionnelle
	reason := "Aucune raison spécifiée"
	if len(parts) > 3 {
		reason = strings.Join(parts[3:], " ")
	}

	err = discord.TimeoutUser(s, m.GuildID, target.ID, duration)
	if err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Timeout non appliqué : %v", err),
		)
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			"Utilisateur %s timeout %v : %s",
			parts[1],
			duration,
			reason,
		),
	)
}
