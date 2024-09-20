package modvoice_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord" // Assurez-vous d'importer votre package pour les vérifications de rôle
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// MoveUserVoiceCommand déplace un utilisateur d'un salon vocal à un autre
func MoveUserVoiceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%smovevoice", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupère le nom du salon vocal cible
	args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(m.Content, command)))
	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier le nom du salon vocal cible.")
		return
	}
	targetVoiceChannelName := strings.Join(args, " ")

	// Récupère le membre qui a envoyé le message
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des informations de l'utilisateur : "+err.Error())
		return
	}

	// Vérifie si le membre est dans un salon vocal
	voiceState, err := s.State.VoiceState(m.GuildID, member.User.ID)
	if err != nil || voiceState == nil {
		s.ChannelMessageSend(m.ChannelID, "Vous n'êtes pas dans un salon vocal.")
		return
	}

	// Recherche le salon vocal cible par nom
	var targetVoiceChannelID string
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des salons : "+err.Error())
		return
	}

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice && channel.Name == targetVoiceChannelName {
			targetVoiceChannelID = channel.ID
			break
		}
	}

	if targetVoiceChannelID == "" {
		s.ChannelMessageSend(m.ChannelID, "Salon vocal cible introuvable.")
		return
	}

	// Déplace l'utilisateur vers le salon vocal cible
	err = s.GuildMemberMove(m.GuildID, member.User.ID, &targetVoiceChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du déplacement de l'utilisateur : "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur a été déplacé vers le salon vocal %s.", targetVoiceChannelName))
}
