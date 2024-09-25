package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Créer un salon
func CreateChannel(s *discordgo.Session, guildID, channelID, channelName string, isVoice bool, duration time.Duration) error {
	channelType := discordgo.ChannelTypeGuildText
	if isVoice {
		channelType = discordgo.ChannelTypeGuildVoice
	}

	channel, err := s.GuildChannelCreate(guildID, channelName, channelType)
	if err != nil {
		return err
	}

	s.ChannelMessageSend(channelID, "Salon créé : <#"+channel.ID+">")

	// Si une durée est définie, supprimer le salon après cette durée
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			s.ChannelDelete(channel.ID)
		}()
	}

	return nil
}

// Supprimer un salon
func DeleteChannel(s *discordgo.Session, guildID, channelID, channelName string) error {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return err
	}

	var channelToDelete *discordgo.Channel
	for _, channel := range channels {
		if channel.Name == channelName {
			channelToDelete = channel
			break
		}
	}

	if channelToDelete == nil {
		return fmt.Errorf("salon non trouvé")
	}

	if _, err := s.ChannelDelete(channelToDelete.ID); err != nil {
		return err
	}

	s.ChannelMessageSend(channelID, "Salon supprimé : "+channelName)
	return nil
}

// Verrouiller ou déverrouiller un salon
func HandleChannelLock(s *discordgo.Session, m *discordgo.MessageCreate, channelID string, duration time.Duration) error {
	// Récupérer le salon directement depuis l'API
	channel, err := s.Channel(channelID)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération du salon (ID: %s) : %v", channelID, err)
	}

	// Récupérer l'ID du rôle @everyone
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération de la guilde : %v", err)
	}
	everyoneRoleID := guild.ID // L'ID du rôle @everyone est l'ID de la guilde

	// Vérifier le statut de verrouillage
	if IsLocked(channel) {
		UnlockChannel(s, channelID)
		// Déverrouiller le salon
	} else {
		// Verrouiller le salon
		LockChannel(s, channelID)
	}

	// Si une durée est définie, déverrouiller après cette durée
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			s.ChannelPermissionSet(channelID, everyoneRoleID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		}()
	}

	return nil
}

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

// GetServerName récupère le nom du serveur à partir de son ID
func GetServerName(s *discordgo.Session, serverID string) (string, error) {
	guild, err := s.Guild(serverID)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du serveur: %w", err)
	}
	return guild.Name, nil
}

// GetChannelName récupère le nom du canal à partir de son ID
func GetChannelName(s *discordgo.Session, channelID string) (string, error) {
	channel, err := s.Channel(channelID)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du canal: %w", err)
	}
	return channel.Name, nil
}
