package channel_commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// CreateTemporaryTextChannelCommand crée un salon textuel temporaire
func CreateTemporaryTextChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%screatetextchannel", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		// Récupère le temps en minutes
		args := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		duration, err := strconv.Atoi(args)
		if err != nil || duration <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un temps valide en minutes.")
			return
		}

		// Crée le salon textuel
		channel, err := s.GuildChannelCreate(m.GuildID, "temp-text-channel", discordgo.ChannelTypeGuildText)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la création du salon : "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Salon textuel temporaire créé : <#"+channel.ID+">")

		// Supprime le salon après la durée spécifiée
		go func() {
			time.Sleep(time.Duration(duration) * time.Minute)
			s.ChannelDelete(channel.ID)
		}()
	}
}
