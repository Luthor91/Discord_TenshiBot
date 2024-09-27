package utility_commands

import (
	"fmt"
	"log"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// PingCommand renvoie le temps de réponse du bot en millisecondes
func PingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sping", config.AppConfig.BotPrefix)
	// Vérifie si le message commence par "!ping"
	if m.Content == command {
		latency := s.HeartbeatLatency().Milliseconds() // Obtenez la latence en millisecondes
		response := fmt.Sprintf("Pong! Latence : %d ms", latency)
		_, err := s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
	}
}
