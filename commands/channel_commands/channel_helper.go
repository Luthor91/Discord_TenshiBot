package channel_commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles si seul ?channel est utilisé
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := `
**Commande ?channel :**

Arguments disponibles :
- **-n [nom]** : Spécifier le nom du salon (par défaut "channel").
- **-v** : Spécifier que c'est un salon vocal.
- **-t [durée]** : Spécifier la durée avant la suppression du salon (exemple : 1h, 30m).
- **-l** : Verrouiller ou déverrouiller le salon.
- **-c** : Créer un nouveau salon.
- **-d** : Supprimer un salon.
`
	s.ChannelMessageSend(channelID, helpMessage)
}

// Récupérer et analyser les arguments de la commande
func parseChannelArgs(m *discordgo.MessageCreate) (string, time.Duration, bool, bool, bool, bool, error) {
	args := strings.Fields(m.Content)
	var (
		duration      time.Duration
		channelName   string
		isVoice       bool
		shouldLock    bool
		createChannel bool
		deleteChannel bool
		err           error
	)

	for i, arg := range args {
		switch arg {
		case "-t":
			if i+1 < len(args) {
				duration, err = utils.ParseDuration(args[i+1])
				if err != nil {
					return "", 0, false, false, false, false, fmt.Errorf("temps non valide")
				}
			}
		case "-n":
			if i+1 < len(args) {
				channelName = args[i+1]
			}
		case "-v":
			isVoice = true
		case "-l":
			shouldLock = true
		case "-c":
			createChannel = true
		case "-d":
			deleteChannel = true
		}
	}

	if channelName == "" {
		channelName = "channel"
	}

	return channelName, duration, isVoice, shouldLock, createChannel, deleteChannel, nil
}

// Créer un salon
func createChannel(s *discordgo.Session, guildID, channelID, channelName string, isVoice bool, duration time.Duration) error {
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
func deleteChannel(s *discordgo.Session, guildID, channelID, channelName string) error {
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
func handleChannelLock(s *discordgo.Session, m *discordgo.MessageCreate, channelID string, duration time.Duration) error {
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
	if discord.IsLocked(channel) {
		// Déverrouiller le salon
		err = s.ChannelPermissionSet(channelID, everyoneRoleID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
		if err != nil {
			return fmt.Errorf("erreur lors du déverrouillage du salon (ID: %s) : %v", channel.ID, err)
		}
		s.ChannelMessageSend(m.ChannelID, "Salon déverrouillé.")
	} else {
		// Verrouiller le salon
		err = s.ChannelPermissionSet(channelID, everyoneRoleID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionSendMessages, 0)
		if err != nil {
			return fmt.Errorf("erreur lors du verrouillage du salon (ID: %s) : %v", channel.ID, err)
		}
		s.ChannelMessageSend(m.ChannelID, "Salon verrouillé.")
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
