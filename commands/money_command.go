package commands

import (
	"fmt"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// MoneyCommand affiche la quantité de monnaie de l'utilisateur
func MoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)

	if m.Content == command {
		money := features.GetUserMoney(m.Author.ID)
		response := fmt.Sprintf("Vous avez %d unités de monnaie.", money)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
