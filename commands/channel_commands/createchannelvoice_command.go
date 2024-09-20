package channel_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// CreateVoiceChannelCommand crée un nouveau salon vocal
func CreateVoiceChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%screatevoicechannel", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		channelName := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if channelName == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nom pour le salon vocal.")
			return
		}

		// Crée le salon vocal
		channel, err := s.GuildChannelCreate(m.GuildID, channelName, discordgo.ChannelTypeGuildVoice)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la création du salon : "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Salon vocal créé : <#"+channel.ID+">")
	}
}
