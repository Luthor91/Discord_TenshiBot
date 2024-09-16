package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// KickCommand expulse un utilisateur mentionné du serveur avec une raison
func KickCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%skick", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Extraire les arguments après la commande (utilisateur à expulser et raison)
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Merci de mentionner l'utilisateur à expulser et de fournir une raison.")
			return
		}

		// Vérifier si l'option verbose "-v" est présente
		verbose := false
		if args[len(args)-1] == "-v" {
			verbose = true
			args = args[:len(args)-1] // Supprimer le "-v" des arguments
		}

		// Extraire l'ID de l'utilisateur mentionné
		userID := args[1]
		if !strings.HasPrefix(userID, "<@") {
			if verbose {
				s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide à expulser.")
			}
			return
		}

		// Nettoyer l'ID de l'utilisateur mentionné
		userID = strings.Trim(userID, "<@!>")

		// Rejoindre la raison
		reason := "Violation des règles"
		if len(args) > 2 {
			reason = strings.Join(args[2:], " ")
		}

		// Expulser l'utilisateur
		err := s.GuildMemberDeleteWithReason(m.GuildID, userID, reason)
		if err != nil {
			log.Printf("Erreur lors de l'expulsion de l'utilisateur : %v", err)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'expulsion de l'utilisateur.")
			}
			return
		}

		// Message de confirmation si verbose est activé
		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été expulsé pour : %s.", userID, reason))
		}
	}
}
