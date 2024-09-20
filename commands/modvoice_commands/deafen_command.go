package modvoice_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord" // Assurez-vous d'importer votre package pour les vérifications de rôle
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// DeafenVoiceCommand permet de déafoner un membre dans un salon vocal
func DeafenVoiceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sdeafenvoice", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupère l'ID de l'utilisateur à déafoner
	userID := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier l'ID de l'utilisateur à déafoner.")
		return
	}

	// Vérifie si l'utilisateur a les permissions nécessaires (optionnel)
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission de déafoner un membre.")
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

	// Déafoner l'utilisateur
	err = s.GuildMemberDeafen(m.GuildID, userID, false)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du déafonnement de l'utilisateur : "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été déafonné.", userID))
}
