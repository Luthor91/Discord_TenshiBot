package moderation_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

// TimeoutCommand est la commande qui applique un timeout
func TimeoutCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%stimeout", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {

		args := strings.Fields(m.Content)
		if len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation: !timeout @mention durée[s/m/h/d]")
			return
		}

		// Vérifier s'il y a une mention d'utilisateur
		if len(m.Mentions) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur.")
			return
		}

		// Récupérer l'ID de l'utilisateur mentionné
		userID := m.Mentions[0].ID
		durationStr := args[2]

		// Parser la durée
		duration, err := utils.ParseDuration(durationStr)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors du parsing de la durée : %v", err))
			return
		}

		// Appliquer le timeout
		err = discord.TimeoutUser(s, m.GuildID, userID, duration)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la mise en timeout de l'utilisateur : %v", err))
			return
		}

		// Informer du succès de la commande
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été mis en timeout pour %v.", userID, duration))
	}
}
