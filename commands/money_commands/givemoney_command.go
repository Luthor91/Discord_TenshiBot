package money_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// GiveMoneyCommand permet d'ajouter une certaine quantité d'argent à un utilisateur mentionné
func GiveMoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?givemoney @mention <quantité>
	prefix := fmt.Sprintf("%sgivemoney", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(m.Mentions) == 0 || len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?givemoney @mention <quantité>")
			return
		}

		// Vérifier que l'utilisateur mentionné n'est pas l'auteur
		mentionedUser := m.Mentions[0]
		if mentionedUser.ID == m.Author.ID {
			s.ChannelMessageSend(m.ChannelID, "Vous ne pouvez pas vous donner de l'argent à vous-même.")
			return
		}

		// Convertir la quantité d'argent en entier
		moneyAmount, err := strconv.Atoi(args[2])
		if err != nil || moneyAmount <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'argent.")
			return
		}

		// Ajouter l'argent à l'utilisateur mentionné
		err = services.NewUserService(controllers.NewUserController()).GiveMoney(m.Author.ID, mentionedUser.ID, moneyAmount)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout de l'argent.")
			return
		}

		// Confirmation à l'utilisateur mentionné
		response := fmt.Sprintf("%s a reçu %d unités de monnaie !", mentionedUser.Mention(), moneyAmount)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
