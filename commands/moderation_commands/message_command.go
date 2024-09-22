package moderation_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/bwmarrin/discordgo"
)

func ModerateMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ne pas réagir à ses propres messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur (adapter cette vérification selon votre logique)
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour exécuter cette commande.")
		return
	}

	// Vérifier si la commande commence par "?msg"
	if !strings.HasPrefix(m.Content, "?msg") {
		return
	}

	// Parsing command
	args := strings.Fields(m.Content)

	var userID string
	var channelID string = m.ChannelID // Par défaut, on prend le salon actuel
	var deleteCount int
	var parseErr error
	var verbose bool // Indicateur pour afficher le message de confirmation

	// Parse command arguments
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-n": // Option pour spécifier un utilisateur
			if i+1 < len(args) {
				userID = args[i+1]
				i++
			}
		case "-c": // Option pour spécifier un salon
			if i+1 < len(args) {
				channelID = args[i+1]
				i++
			}
		case "-d": // Option pour spécifier le nombre de messages à supprimer
			if i+1 < len(args) {
				deleteCount, parseErr = strconv.Atoi(args[i+1])
				if parseErr != nil || deleteCount <= 0 {
					s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre valide de messages à supprimer.")
					return
				}
				i++
			}
		case "-v": // Option pour afficher un message de confirmation
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

	// Résoudre l'ID du salon si seul le nom est fourni
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
		// Si userID est vide, supprimer tous les messages sinon filtrer par l'utilisateur spécifié
		if userID == "" || message.Author.ID == userID {
			err := s.ChannelMessageDelete(channelID, message.ID)
			if err != nil {
				fmt.Printf("Erreur lors de la suppression du message : %v\n", err)
			}
		}
	}

	// Afficher le message de confirmation si l'option "-v" est présente
	if verbose {
		// Récupérer les informations du canal
		channel, err := s.Channel(channelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Messages supprimés, mais erreur lors de la récupération des informations du salon.")
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Messages supprimés avec succès dans le salon %s.", channel.Name))
		}
	}
}
