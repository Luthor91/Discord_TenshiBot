package modvoice_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// KickVoiceCommand permet de kick un utilisateur d'un salon vocal
func KickVoiceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%skickvoice", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupère l'identifiant de l'utilisateur à kick
	userID := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
	if userID == "" {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier l'ID de l'utilisateur à kick.")
		return
	}

	// Vérifie si l'utilisateur est dans un salon vocal
	_, err = s.GuildMember(m.GuildID, userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des informations de l'utilisateur : "+err.Error())
		return
	}

	// Déplace l'utilisateur vers un salon "null" (le déconnecte)
	if err := s.GuildMemberMove(m.GuildID, userID, nil); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du kick de l'utilisateur : "+err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été kické du salon vocal.", userID))
}
