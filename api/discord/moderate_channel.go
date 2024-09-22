package discord

import "github.com/bwmarrin/discordgo"

func LockChannel(s *discordgo.Session, channelID string) error {
	// Changer les permissions pour verrouiller le salon
	return s.ChannelPermissionSet(channelID, "@everyone", discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
}

func UnlockChannel(s *discordgo.Session, channelID string) error {
	// Réinitialiser les permissions pour déverrouiller le salon
	return s.ChannelPermissionSet(channelID, "@everyone", discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
}

func IsLocked(channel *discordgo.Channel) bool {
	for _, perm := range channel.PermissionOverwrites {
		if perm.ID == "@everyone" && perm.Deny&discordgo.PermissionSendMessages > 0 {
			return true
		}
	}
	return false
}
