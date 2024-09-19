package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// DeleteCommand supprime un nombre spécifié de messages dans le salon où la commande est exécutée
func DeleteCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sdelete", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Extraire les arguments après la commande (nombre de messages à supprimer)
		args := strings.Fields(m.Content)

		// Vérifier si l'argument du nombre de messages est présent
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Merci de spécifier le nombre de messages à supprimer.")
			return
		}

		// Vérifier si l'option verbose "-v" est présente
		verbose := false
		if args[len(args)-1] == "-v" {
			verbose = true
			args = args[:len(args)-1] // Supprimer le "-v" des arguments
		}

		// Vérifier que le premier argument est un nombre
		numMessagesStr := args[1]
		numMessages, err := strconv.Atoi(numMessagesStr)
		if err != nil || numMessages <= 0 {
			if verbose {
				s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un nombre valide de messages à supprimer (doit être supérieur à 0).")
			}
			return
		}

		// Obtenir les messages récents dans le canal
		messages, err := s.ChannelMessages(m.ChannelID, numMessages+1, "", "", "")
		if err != nil {
			log.Printf("Erreur lors de la récupération des messages: %v", err)
			return
		}

		// Collecter les IDs des messages à supprimer
		var messageIDs []string
		for _, msg := range messages {
			messageIDs = append(messageIDs, msg.ID)
		}

		// Supprimer les messages
		if len(messageIDs) > 1 {
			err = s.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
			if err != nil {
				log.Printf("Erreur lors de la suppression des messages: %v", err)
				if verbose {
					s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression des messages.")
				}
				return
			}
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Supprimé %d messages.", numMessages))
			}
		} else if len(messageIDs) == 1 {
			err = s.ChannelMessageDelete(m.ChannelID, messageIDs[0])
			if err != nil {
				log.Printf("Erreur lors de la suppression du message: %v", err)
				if verbose {
					s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du message.")
				}
			} else if verbose {
				s.ChannelMessageSend(m.ChannelID, "Message supprimé.")
			}
		}
	}
}
