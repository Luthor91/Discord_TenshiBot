package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// timeoutUser met un utilisateur en timeout pour une durée donnée
func TimeoutUser(s *discordgo.Session, guildID, userID string, duration time.Duration) error {
	timeoutUntil := time.Now().Add(duration)
	err := s.GuildMemberTimeout(guildID, userID, &timeoutUntil)
	if err != nil {
		return err
	}
	return nil
}
