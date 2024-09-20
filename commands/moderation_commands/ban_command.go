package moderation_commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// BanCommand bannit un utilisateur mentionné du serveur avec une raison
func BanCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sban", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Extraire les arguments après la commande (utilisateur à bannir et raison)
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Merci de mentionner l'utilisateur à bannir et de fournir une raison.")
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
				s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide à bannir.")
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

		// Bannir l'utilisateur
		err := s.GuildBanCreateWithReason(m.GuildID, userID, reason, 0)
		if err != nil {
			log.Printf("Erreur lors du bannissement de l'utilisateur : %v", err)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, "Erreur lors du bannissement de l'utilisateur.")
			}
			return
		}

		// Message de confirmation si verbose est activé
		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été banni pour : %s.", userID, reason))
		}
	}
}
