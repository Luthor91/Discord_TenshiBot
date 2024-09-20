package stat_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// BotStatsCommand affiche les statistiques du bot
func BotStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sbotstat", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Compter le nombre total de membres
	totalUsers := 0
	for _, guild := range s.State.Guilds {
		guild, err := s.Guild(guild.ID) // Optionnel: récupérer la guilde actuelle pour assurer une synchronisation
		if err == nil {
			totalUsers += guild.MemberCount
		}
	}

	// Exemples de statistiques
	response := fmt.Sprintf("**Statistiques du bot :**\n- Serveurs: %d\n- Utilisateurs: %d",
		len(s.State.Guilds), totalUsers)
	s.ChannelMessageSend(m.ChannelID, response)
}
