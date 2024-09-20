package channel_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// DeleteChannelByNameCommand supprime un salon écrit ou vocal par son nom
func DeleteChannelByNameCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sdeletechannel", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		channelName := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if channelName == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier le nom du salon à supprimer.")
			return
		}

		channels, err := s.GuildChannels(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des salons : "+err.Error())
			return
		}

		var channelToDelete *discordgo.Channel
		for _, channel := range channels {
			if channel.Name == channelName {
				channelToDelete = channel
				break
			}
		}

		if channelToDelete == nil {
			s.ChannelMessageSend(m.ChannelID, "Salon non trouvé.")
			return
		}

		// Supprime le salon
		if _, err := s.ChannelDelete(channelToDelete.ID); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du salon : "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Salon supprimé : "+channelName)
	}
}
