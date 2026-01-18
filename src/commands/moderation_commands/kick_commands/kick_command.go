package kick_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// KickCommand gère la commande ?kick
// Usage : ?kick <user> [raison]
func KickCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "kick"

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// ?kick
	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		showHelpMessage(s, m.ChannelID)
		return
	}

	target := discord.HandleTarget(s, m, parts[1])
	if target == nil {
		return
	}

	// Raison optionnelle
	reason := "Aucune raison spécifiée"
	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}

	discord.KickUser(s, m, target.ID, reason)

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Utilisateur %s expulsé : %s", parts[1], reason),
	)
}
