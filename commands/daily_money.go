package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

// DailyMoneyCommand exécute la commande pour donner une récompense quotidienne
func DailyMoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sdaily", config.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		canReceive, timeLeft := features.CanReceiveDailyReward(m.Author.ID)

		if !canReceive {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous devez attendre encore %v avant de réclamer votre prochaine récompense quotidienne.", m.Author.Username, timeLeft.Round(time.Minute)))
			return
		}

		// Générer un montant aléatoire entre 10 et 100
		randomAmount := rand.Intn(91) + 10

		// Donner la récompense à l'utilisateur
		features.GiveDailyMoney(m.Author.ID, randomAmount)

		// Envoyer un message de confirmation
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous avez reçu %d pièces aujourd'hui !", m.Author.Username, randomAmount))
	}
}
