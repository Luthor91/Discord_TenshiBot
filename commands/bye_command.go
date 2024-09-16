package commands

import (
	"fmt"
	"log"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// ByeCommand répond "Good Bye👋" lorsque l'utilisateur tape "!bye"
func ByeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sbye", config.AppConfig.BotPrefix)

	if m.Content == command {
		_, err := s.ChannelMessageSend(m.ChannelID, "Good Bye👋")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
	}
}
