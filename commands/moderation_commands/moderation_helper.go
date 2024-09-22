package moderation_commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// ParseArgs parses the command line arguments
func parseArgs(args []string) (map[string]string, error) {
	parsedArgs := make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n", "-b", "-w", "-k", "-m", "-d", "-to", "-mv":
			if i+1 < len(args) {
				parsedArgs[arg] = args[i+1]
				i++
			} else {
				return nil, errors.New("missing value for " + arg)
			}
		case "-t":
			if i+1 < len(args) {
				parsedArgs[arg] = args[i+1]
				i++
			} else {
				return nil, errors.New("missing duration value for -t")
			}
		case "-r", "-rw":
			parsedArgs[arg] = "true" // Just a flag, doesn't need a value
		}
	}
	return parsedArgs, nil
}

// Helper functions for action handling
func banUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	err := s.GuildBanCreateWithReason(m.GuildID, userID, reason, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du bannissement de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s banni avec succès.", userID))
	}
}

func warnUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissement envoyé à l'utilisateur %s : %s", userID, reason))
}

func kickUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, reason string) {
	err := s.GuildMemberDeleteWithReason(m.GuildID, userID, reason)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du kick de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s kické avec succès.", userID))
	}
}

func muteUser(s *discordgo.Session, m *discordgo.MessageCreate, userID string, duration time.Duration, reason string) {
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

func deafenUser(s *discordgo.Session, m *discordgo.MessageCreate, userID string, duration time.Duration, reason string) {
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

func timeoutUser(s *discordgo.Session, m *discordgo.MessageCreate, userID string, duration time.Duration, reason string) {
	// Implement timeout logic here (DiscordGo does not natively support timeout yet)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Timeout de l'utilisateur %s pour %s : %s", userID, duration.String(), reason))
}

func moveUser(s *discordgo.Session, m *discordgo.MessageCreate, userID, targetChannel string) {
	err := s.GuildMemberMove(m.GuildID, userID, &targetChannel)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du déplacement de l'utilisateur.")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Utilisateur %s déplacé avec succès.", userID))
	}
}

func resetAllUserStatus(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	// Implement reset logic to clear all punishments (ban, mute, deafen, etc.)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Statut de l'utilisateur %s réinitialisé.", userID))
}

func resetUserWarnings(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	// Implement reset warnings logic
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissements de l'utilisateur %s réinitialisés.", userID))
}
