package moderation_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// ModerateMessageCommand gère les commandes de modération des messages
func ModerateMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour exécuter cette commande.")
		return
	}

	command := fmt.Sprintf("%smessage", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%smsg", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	// Extraction des arguments
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var userID, channelID string = m.ChannelID, m.ChannelID
	var deleteCount int
	var verbose bool

	// Analyse des arguments extraits
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			userID = discord.HandleTarget(s, m, arg.Value).ID
		case "-c":
			channelID = arg.Value
		case "-d":
			deleteCount, err = strconv.Atoi(arg.Value)
			if err != nil || deleteCount <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre valide de messages à supprimer.")
				return
			}
		case "-v":
			verbose = true
		default:
			// Ignorer les arguments non reconnus
		}
	}

	// Validation de l'input
	if deleteCount == 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre de messages à supprimer avec -d.")
		return
	}

	// Résoudre l'ID du salon si seulement le nom est fourni
	if channelID != "" && !strings.HasPrefix(channelID, "<#") {
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

	// Récupérer les messages à supprimer
	messages, err := s.ChannelMessages(channelID, deleteCount, "", "", "")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des messages.")
		return
	}

	// Supprimer les messages selon l'utilisateur ou le canal spécifié
	for _, message := range messages {
		if userID == "" || message.Author.ID == userID {
			err := s.ChannelMessageDelete(channelID, message.ID)
			if err != nil {
				fmt.Printf("Erreur lors de la suppression du message : %v\n", err)
			}
		}
	}

	// Afficher le message de confirmation si l'option "-v" est présente
	if verbose {
		channel, err := s.Channel(channelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Messages supprimés, mais erreur lors de la récupération des informations du salon.")
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Messages supprimés avec succès dans le salon %s.", channel.Name))
		}
	}
}
