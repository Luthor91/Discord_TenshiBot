package game_commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

func CollectInvestCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%scollect", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer les informations de l'utilisateur
	userService := services.NewUserService()
	user, err := userService.GetUserByDiscordID(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la récupération des informations utilisateur : %v", err))
		return
	}

	// Récupérer le dernier investissement de l'utilisateur
	investmentService := services.NewInvestmentService()
	lastInvestment, err := investmentService.GetLastInvestmentUser(user.UserDiscordID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Aucun investissement trouvé.")
		return
	}

	// Calculer le temps écoulé depuis l'investissement
	hoursSinceInvestment := time.Since(lastInvestment.InvestedAt).Hours()

	// S'assurer que le temps ne dépasse pas 6 heures
	if hoursSinceInvestment > 6 {
		hoursSinceInvestment = 6
	}

	// Calculer le multiplicateur basé sur le temps écoulé
	// Par exemple : 1.0 à 0 heures, 2.0 à 6 heures
	multiplier := 1 + (hoursSinceInvestment / 6) // Le multiplicateur variera entre 1 et 2

	// Calculer le montant de retour
	returnAmount := int(float64(lastInvestment.Amount) * multiplier)

	// Mettre à jour le solde de l'utilisateur
	newMoneyAmount := user.Money + returnAmount
	err = userService.AddMoney(m.Author.ID, newMoneyAmount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la mise à jour du solde : %v", err))
		return
	}

	// Supprimer l'investissement après récupération
	err = investmentService.DeleteInvestment(lastInvestment.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression de l'investissement.")
		return
	}

	// Envoyer un message de confirmation à l'utilisateur
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez récupéré %d avec un facteur de %.2f !", returnAmount, multiplier))
}
