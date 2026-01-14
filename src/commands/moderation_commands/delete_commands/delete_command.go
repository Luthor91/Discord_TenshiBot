package delete_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// DeleteMessageCommand gère les commandes de suppression des messages
func DeleteMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sdelete", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%sdel", config.AppConfig.BotPrefix)

	var prefix string
	if strings.HasPrefix(m.Content, command) {
		prefix = command
	} else if strings.HasPrefix(m.Content, commandAlias) {
		prefix = commandAlias
	} else {
		return
	}

	// Extraction des arguments
	parsedArgs, err := discord.ExtractArguments(m.Content, prefix)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}


	// Initialiser avec le salon actuel comme valeur par défaut
	var userID string = ""
	var channelID string = m.ChannelID // Utiliser le salon actuel par défaut
	var deleteCount int
	var verbose bool

	// Analyse des arguments extraits
	for _, arg := range parsedArgs {
		switch arg.Arg {
			case "-v":
				verbose = true
			case "-n":
				if arg.Value != "" {
					userID, err = resolveUserID(s, m.GuildID, arg.Value)
					if err != nil && verbose {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la résolution du pseudo : %v", err))
						return
					}
				}
			case "-c":
				if arg.Value != "" {
					tempChannelID, err := resolveChannelID(s, m.GuildID, arg.Value)
					if err != nil && verbose {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la résolution du salon : %v", err))
						return
					}
					channelID = tempChannelID
				}
			case "-d":
				deleteCount, err = strconv.Atoi(arg.Value)
				if err != nil || deleteCount <= 0 && verbose {
					s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre valide de messages à supprimer.")
					return
				}
			default:
				// Ignorer les arguments non reconnus
		}
	}

	// Validation de l'input
	if deleteCount == 0 && verbose {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre de messages à supprimer avec -d.")
		return
	}

	// Récupérer les messages à supprimer
	messages, err := s.ChannelMessages(channelID, deleteCount, "", "", "")
	if err != nil && verbose{
		errMsg := fmt.Sprintf("Erreur lors de la récupération des messages: %v", err)
		s.ChannelMessageSend(m.ChannelID, errMsg)
		return
	}

	var deletedCount int = 0
	var messageIDs []string

	for _, message := range messages {
		if userID == "" || message.Author.ID == userID {
			messageIDs = append(messageIDs, message.ID)
			deletedCount++
		}
	}

	if len(messageIDs) > 0 {
		if len(messageIDs) > 1 {
			err = s.ChannelMessagesBulkDelete(channelID, messageIDs)
			if err != nil {

				for _, msgID := range messageIDs {
					err := s.ChannelMessageDelete(channelID, msgID)
					if err != nil && verbose {
						fmt.Printf("Erreur lors de la suppression du message %s: %v\n", msgID, err)
					}
				}
			}
		} else if len(messageIDs) == 1 {
			err := s.ChannelMessageDelete(channelID, messageIDs[0])
			if err != nil && verbose {
				fmt.Printf("Erreur lors de la suppression du message: %v\n", err)
			}
		}
	}

}


