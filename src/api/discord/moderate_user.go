package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// timeoutUser met un utilisateur en timeout pour une durée donnée
func TimeoutUser(s *discordgo.Session, guildID, userID string, duration time.Duration) error {
	timeoutUntil := time.Now().Add(duration)
	err := s.GuildMemberTimeout(guildID, userID, &timeoutUntil)
	if err != nil {
		return err
	}
	return nil
}

// Helper functions for action handling
func BanUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	err := s.GuildBanCreateWithReason(m.GuildID, userID, reason, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du bannissement de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s banni avec succès.", userID))
	}
}

func KickUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	err := s.GuildMemberDeleteWithReason(m.GuildID, userID, reason)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du kick de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s kické avec succès.", userID))
	}
}

func MuteUser(s *discordgo.Session, m *discordgo.MessageCreate, userID string, duration time.Duration, reason string) {
	err := s.GuildMemberMute(m.GuildID, userID, true)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du mute de l'utilisateur.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s mute pour %s : %s", userID, duration.String(), reason))
	// Optionally, schedule unmute after duration
	time.AfterFunc(duration, func() {
		s.GuildMemberMute(m.GuildID, userID, false)
	})
}

func DeafenUser(s *discordgo.Session, m *discordgo.MessageCreate, userID string, duration time.Duration, reason string) {
	err := s.GuildMemberDeafen(m.GuildID, userID, true)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du deafen de l'utilisateur.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s deafen pour %s : %s", userID, duration.String(), reason))
	// Optionally, schedule undeafen after duration
	time.AfterFunc(duration, func() {
		s.GuildMemberDeafen(m.GuildID, userID, false)
	})
}

func MoveUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, targetChannel string) {
	err := s.GuildMemberMove(m.GuildID, userID, &targetChannel)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du déplacement de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s déplacé avec succès.", userID))
	}
}

// FindUserByUsername recherche un utilisateur dans un serveur Discord par son pseudo
func FindUserByUsername(s *discordgo.Session, guildID, username string) (*discordgo.User, error) {
	guild, err := s.Guild(guildID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du serveur : %w", err)
	}

	for _, member := range guild.Members {
		if strings.EqualFold(member.User.Username, username) {
			return member.User, nil
		}
	}
	return nil, fmt.Errorf("utilisateur `%s` non trouvé dans le serveur", username)
}
