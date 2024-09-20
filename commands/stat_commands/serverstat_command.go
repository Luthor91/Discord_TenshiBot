package stat_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// ServerStatsCommand affiche des statistiques sur le serveur
func ServerStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sserverstat", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des statistiques du serveur.")
		return
	}

	response := fmt.Sprintf("**Statistiques du serveur :**\n- Membres: %d\n- Rôles: %d\n- Canaux: %d",
		len(guild.Members), len(guild.Roles), len(guild.Channels))
	s.ChannelMessageSend(m.ChannelID, response)
}
