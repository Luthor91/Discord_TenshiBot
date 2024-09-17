package commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// AddBadWordCommand ajoute un mot dans la liste des "badwords"
func AddBadWordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%saddbadword", config.AppConfig.BotPrefix)
	// Vérifie si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Récupère le mot à ajouter
		word := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if word == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à ajouter.")
			return
		}

		// Ajoute le mot dans la liste des "badwords"
		features.AddBadWord(word)
	}
}
