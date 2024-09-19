package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// DailyMoneyCommand exécute la commande pour donner une récompense quotidienne
func DailyMoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sdaily", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Assurez-vous d'initialiser le générateur de nombres aléatoires
		rand.Seed(time.Now().UnixNano())

		// Vérifier si l'utilisateur peut recevoir la récompense quotidienne
		canReceive, timeLeft, err := services.CanReceiveDailyReward(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la vérification de la récompense quotidienne : "+err.Error())
			return
		}

		if !canReceive {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous devez attendre encore %v avant de réclamer votre prochaine récompense quotidienne.", m.Author.Username, timeLeft.Round(time.Minute)))
			return
		}

		// Générer un montant aléatoire entre 10 et 100
		randomAmount := rand.Intn(91) + 10

		// Donner la récompense à l'utilisateur
		if err := services.GiveDailyMoney(m.Author.ID, randomAmount); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'attribution de la récompense : "+err.Error())
			return
		}

		// Envoyer un message de confirmation
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous avez reçu %d pièces aujourd'hui !", m.Author.Username, randomAmount))
	}
}
