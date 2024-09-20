package experience_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// SetXPCommand permet de définir une quantité spécifique d'XP pour l'auteur
func SetXPCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier les droits d'administrateur
	isAdmin, err := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isAdmin {
		return
	}

	// Commande de base : ?setxp <quantité>
	prefix := fmt.Sprintf("%ssetxp", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?setxp <quantité>")
			return
		}

		// Convertir la quantité d'XP en entier
		xpAmount, err := strconv.Atoi(args[1])
		if err != nil || xpAmount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'XP.")
			return
		}

		// Définir l'XP pour l'utilisateur
		err = services.NewUserService(controllers.NewUserController()).SetExperience(m.Author.ID, xpAmount)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la définition de l'XP.")
			return
		}

		// Confirmation à l'utilisateur
		response := fmt.Sprintf("Votre XP a été défini à %d.", xpAmount)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
