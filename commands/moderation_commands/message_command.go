package moderation_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ModerateMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Parsing command
	args := strings.Fields(m.Content)

	var userID string
	var channelID string
	var deleteCount int
	var err error

	// Parse command arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n":
			if i+1 < len(args) {
				userID = args[i+1]
				i++
			}
		case "-c":
			if i+1 < len(args) {
				channelID = args[i+1]
				i++
			}
		case "-d":
			if i+1 < len(args) {
				deleteCount, err = strconv.Atoi(args[i+1])
				if err != nil || deleteCount <= 0 {
					s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre valide de messages à supprimer.")
					return
				}
				i++
			}
		default:
			// Ignore unrecognized arguments
		}
	}

	// Validate input
	if deleteCount == 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre de messages à supprimer avec -d.")
		return
	}
	if userID == "" && channelID == "" {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier soit un utilisateur avec -n, soit un canal avec -c.")
		return
	}

	// Resolve channel ID if only name is provided
	if channelID != "" && !strings.HasPrefix(channelID, "<#") {
		// Search for the channel by name
		channels, err := s.GuildChannels(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des salons.")
			return
		}
		for _, channel := range channels {
			if channel.Name == channelID {
				channelID = channel.ID
				break
			}
		}
	}

	// Default to current channel if channelID is still empty
	if channelID == "" {
		channelID = m.ChannelID
	}

	// Fetch messages to delete
	messages, err := s.ChannelMessages(channelID, deleteCount, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des messages.")
		return
	}

	// Delete messages based on user or channel specification
	for _, message := range messages {
		if userID == "" || message.Author.ID == userID {
			err := s.ChannelMessageDelete(channelID, message.ID)
			if err != nil {
				fmt.Printf("Erreur lors de la suppression du message : %v\n", err)
			}
		}
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Messages supprimés avec succès dans le canal %s.", channelID))
}
