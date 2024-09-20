package money_commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
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
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Créer un contrôleur d'utilisateur
		userController := controllers.NewUserController()

		// Récupérer les informations de l'utilisateur
		user, err := userController.GetUserByDiscordID(m.Author.ID) // Ajoute cette méthode dans ton UserController
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de vos informations : "+err.Error())
			return
		}

		// Vérifier si l'utilisateur peut recevoir la récompense quotidienne
		canReceive, timeLeft := services.NewUserService(userController).CanReceiveDailyReward(user)

		if !canReceive {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous devez attendre encore %v avant de réclamer votre prochaine récompense quotidienne.", m.Author.Username, timeLeft.Round(time.Minute)))
			return
		}

		// Générer un montant aléatoire entre 10 et 100
		randomAmount := rng.Intn(91) + 10

		// Donner la récompense à l'utilisateur
		if err := services.NewUserService(userController).UpdateDailyMoney(user, randomAmount); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'attribution de la récompense : "+err.Error())
			return
		}

		// Envoyer un message de confirmation
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous avez reçu %d pièces aujourd'hui !", m.Author.Username, randomAmount))
	}
}
