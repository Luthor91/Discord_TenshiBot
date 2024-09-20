package channel_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// DeLockChannelCommand déverrouille un salon
func DeLockChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sdelockchan", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}
	// Vérifie si un salon est mentionné ou utilise le salon actuel
	channelID := m.ChannelID
	if len(m.Mentions) > 0 {
		channelID = m.Mentions[0].ID
	}

	// Modifie les permissions pour déverrouiller le salon (à adapter selon les permissions souhaitées)
	s.ChannelMessageEdit(m.ChannelID, channelID, "Le salon est déverrouillé.")
}
