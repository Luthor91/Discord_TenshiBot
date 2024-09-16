package commands

import (
	"fmt"
	"log"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// PingCommand envoie "Pong!" quand l'utilisateur tape "!ping"
func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sping", config.AppConfig.BotPrefix)
	// Vérifie si le message commence par "!ping"
	if m.Content == command {
		_, err := s.ChannelMessageSend(m.ChannelID, "Pong!")
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}
