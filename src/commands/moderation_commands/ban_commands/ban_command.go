package ban_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// BanCommand gère la commande ?ban
// Usage : ?ban <user> [raison]
func BanCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "ban"

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// ?ban
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

	discord.BanUser(s, m, target.ID, reason)

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf("Utilisateur %s banni : %s", target.Username, reason),
	)
}

