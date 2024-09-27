package daily_commands

import (
	"fmt"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// MoneyCommand gère la commande quotidienne
func DailyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := config.AppConfig.BotPrefix
	command := fmt.Sprintf("%sdaily", prefix)

	// Vérifie si le message commence par la commande
	if m.Content == command {
		handleDailyReward(s, m) // Appelle la fonction associée pour la récompense quotidienne
	}
}
