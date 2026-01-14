package move_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// MoveCommand gère la commande ?move
// Usage : ?move <user> <salon>
func MoveCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "move"

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// ?move
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

	channel, err := discord.HandleChannel(s, m, parts[2], discord.VoiceChannel)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	discord.MoveUser(s, m, target.ID, channel.ID)

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			"Utilisateur %s déplacé vers %s",
			parts[1],
			channel.Name,
		),
	)
}

