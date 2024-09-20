package modvoice_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// MuteVoiceCommand permet de mute un membre dans un salon vocal
func MuteVoiceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%smute", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupère l'ID de l'utilisateur à mute
	userID := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier l'ID de l'utilisateur à mute.")
		return
	}

	// Vérifie si l'utilisateur est dans un salon vocal
	_, err = s.GuildMember(m.GuildID, userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des informations de l'utilisateur : "+err.Error())
		return
	}

	// Vérifie si le membre est dans un salon vocal
	voiceState, err := s.State.VoiceState(m.GuildID, userID)
	if err != nil || voiceState == nil {
		s.ChannelMessageSend(m.ChannelID, "L'utilisateur n'est pas dans un salon vocal.")
		return
	}

	// Muter l'utilisateur
	err = s.GuildMemberMute(m.GuildID, userID, true)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du mute de l'utilisateur : "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été mute.", userID))
}
