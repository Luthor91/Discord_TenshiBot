package commands

import (
	"fmt"
	"log"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// HelpCommand répond "Hello World😃" lorsque l'utilisateur tape "!help"
func HelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%shelp", config.AppConfig.BotPrefix)

	if m.Content == command {
		_, err := s.ChannelMessageSend(m.ChannelID, "Hello World😃")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
	}
}
