package moderation_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// ResetWarnCommand réinitialise les avertissements d'un utilisateur
func ResetWarnCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorer les messages du bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if !isMod {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sresetwarn", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer la mention
	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?resetwarn <@mention_user>")
		return
	}

	// Récupérer l'utilisateur mentionné
	mention := parts[1]
	if !strings.HasPrefix(mention, "<@") || !strings.HasSuffix(mention, ">") {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide.")
		return
	}
	mentionID := mention[2 : len(mention)-1] // Extrait l'ID de la mention

	// Créer une instance de WarnController
	warnController := controllers.NewWarnController()

	// Créer une instance de WarnService
	warnService := services.NewWarnService(warnController, s, m.GuildID)

	// Réinitialiser les avertissements pour l'utilisateur
	err := warnService.ResetWarns(mentionID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la réinitialisation des avertissements.")
		return
	}

	// Envoyer un message confirmant la réinitialisation
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Tous les avertissements pour <@%s> ont été réinitialisés.", mentionID))
}
