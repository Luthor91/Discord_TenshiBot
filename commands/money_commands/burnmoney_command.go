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

// BurnMoneyCommand retire un montant d'argent à l'utilisateur
func BurnMoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?burnmoney <montant>
	prefix := fmt.Sprintf("%sburnmoney", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?burnmoney <montant>")
			return
		}

		// Convertir le montant en entier
		moneyAmount, err := strconv.Atoi(args[1])
		if err != nil || moneyAmount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide d'argent.")
			return
		}

		// Récupérer l'argent de l'utilisateur
		amount, err := services.NewUserService(controllers.NewUserController()).GetMoney(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé.")
			return
		}

		// Vérifier si l'utilisateur a assez d'argent à retirer
		if amount < moneyAmount {
			s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour retirer ce montant.")
			return
		}

		// Définir l'argent pour l'utilisateur
		newMoney := amount - moneyAmount
		err = services.NewUserService(controllers.NewUserController()).SetMoney(m.Author.ID, newMoney)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'argent.")
			return
		}

		// Confirmation à l'utilisateur
		response := fmt.Sprintf("Vous avez retiré %d. Votre nouvel argent est de %d.", moneyAmount, newMoney)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
