package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

// ReminderCommand définit un rappel après une durée donnée, en extrayant la durée et le message du contenu du message
// Prend en charge les suffixes: s (secondes), m (minutes), h (heures), d (jours)
func ReminderCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Extraire le préfixe et la commande
	command := fmt.Sprintf("%sreminder", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les paramètres de la commande
	// Format attendu : !reminder <durée><s|m|h|d> <message>
	args := strings.TrimSpace(m.Content[len(command):])
	parts := strings.SplitN(args, " ", 2)

	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Format incorrect. Utilisez: !reminder <durée><s|m|h|d> <message>")
		return
	}

	// Extraire la durée et le message
	durationStr := parts[0]
	message := parts[1]

	// Convertir la durée basée sur le suffixe (s, m, h, d)
	duration, err := utils.ParseDuration(durationStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Durée invalide. Utilisez une durée avec un suffixe : s (secondes), m (minutes), h (heures), d (jours).")
		return
	}

	// Définir le rappel
	go func() {
		time.Sleep(duration)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Rappel pour %s : %s", m.Author.Mention(), message))
	}()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Rappel défini pour %s après %s.", m.Author.Mention(), duration.String()))
}
