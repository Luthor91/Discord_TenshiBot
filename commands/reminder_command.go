package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Reminder définit un rappel après une durée donnée
func ReminderCommand(s *discordgo.Session, m *discordgo.MessageCreate, duration time.Duration, message string) {
	go func() {
		time.Sleep(duration)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Rappel pour %s : %s", m.Author.Mention(), message))
	}()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Rappel défini pour %s après %s.", m.Author.Mention(), duration.String()))
}
