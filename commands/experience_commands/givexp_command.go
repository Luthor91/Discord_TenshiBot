package experience_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// GiveXPCommand permet d'ajouter une certaine quantité d'XP à un utilisateur mentionné
func GiveXPCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?givexp @mention <quantité>
	prefix := fmt.Sprintf("%sgivexp", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(m.Mentions) == 0 || len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?givexp @mention <quantité>")
			return
		}

		// Vérifier que l'utilisateur mentionné n'est pas l'auteur
		mentionedUser := m.Mentions[0]
		if mentionedUser.ID == m.Author.ID {
			s.ChannelMessageSend(m.ChannelID, "Vous ne pouvez pas vous donner de l'XP à vous-même.")
			return
		}

		// Convertir la quantité d'XP en entier
		xpAmount, err := strconv.Atoi(args[2])
		if err != nil || xpAmount <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'XP.")
			return
		}

		// Ajouter l'XP à l'utilisateur mentionné
		err = services.NewUserService(controllers.NewUserController()).GiveXP(m.Author.ID, mentionedUser.ID, xpAmount)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout de l'XP.")
			return
		}

		// Confirmation à l'utilisateur mentionné
		response := fmt.Sprintf("%s a reçu %d XP !", mentionedUser.Mention(), xpAmount)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
