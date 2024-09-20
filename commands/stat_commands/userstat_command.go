package stat_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// UserStatsCommand affiche des statistiques sur un utilisateur
func UserStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%suserstat", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}
	member, err := s.GuildMember(m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des statistiques de l'utilisateur.")
		return
	}

	response := fmt.Sprintf("**Statistiques de l'utilisateur :**\n- Nom: %s\n- Rôles: %d",
		member.User.Username, len(member.Roles))
	s.ChannelMessageSend(m.ChannelID, response)
}
