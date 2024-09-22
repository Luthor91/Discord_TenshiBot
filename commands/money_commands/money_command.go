package money_commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

func MoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parser la commande avec des flags
	prefix := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)

		// Initialisation des variables
		var targetUser *discordgo.User
		var err error
		var amount int

		// Récupérer les arguments de la commande
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "-n":
				// Récupérer le nom ou la mention de la target
				if len(m.Mentions) > 0 {
					targetUser = m.Mentions[0]
				} else {
					s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur avec -n.")
					return
				}
			case "-r":
				// Retirer de l'argent à l'utilisateur actuel
				if len(args) > i+1 {
					amount, err = strconv.Atoi(args[i+1])
					if err != nil || amount <= 0 {
						s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour burn.")
						return
					}
					userMoney, err := services.NewUserService(controllers.NewUserController()).GetMoney(m.Author.ID)
					if err != nil || userMoney < amount {
						s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour retirer ce montant.")
						return
					}
					newMoney := userMoney - amount
					services.NewUserService(controllers.NewUserController()).SetMoney(m.Author.ID, newMoney)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez retiré %d unités. Nouveau solde : %d", amount, newMoney))
					return
				}
			case "-d":
				// Récompense quotidienne
				userController := controllers.NewUserController()
				user, err := userController.GetUserByDiscordID(m.Author.ID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de vos informations.")
					return
				}
				canReceive, timeLeft := services.NewUserService(userController).CanReceiveDailyReward(user)
				if !canReceive {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous devez attendre encore %v avant de réclamer la prochaine récompense quotidienne.", timeLeft.Round(time.Minute)))
					return
				}
				randomAmount := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(91) + 10
				services.NewUserService(userController).UpdateDailyMoney(user, randomAmount)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez reçu %d unités aujourd'hui !", randomAmount))
				return
			case "-m":
				// Afficher son propre argent
				money, err := services.NewUserService(controllers.NewUserController()).GetMoney(m.Author.ID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de votre solde.")
					return
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez %d unités de monnaie.", money))
				return
			case "-g":
				// Donner de l'argent à un utilisateur mentionné
				if targetUser == nil {
					s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur avec -n.")
					return
				}
				if len(args) > i+1 {
					amount, err = strconv.Atoi(args[i+1])
					if err != nil || amount <= 0 {
						s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour donner.")
						return
					}
					services.NewUserService(controllers.NewUserController()).GiveMoney(m.Author.ID, targetUser.ID, amount)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez donné %d unités à %s.", amount, targetUser.Username))
					return
				}
			case "-s":
				// Vérifier si l'utilisateur est admin avant de permettre la commande -s
				isAdmin, err := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
				if err != nil || !isAdmin {
					s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission de définir la monnaie.")
					return
				}
				// Définir l'argent d'un utilisateur (réservé aux admins)
				if targetUser == nil {
					s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur avec -n.")
					return
				}
				if len(args) > i+1 {
					amount, err = strconv.Atoi(args[i+1])
					if err != nil || amount < 0 {
						s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.")
						return
					}
					services.NewUserService(controllers.NewUserController()).SetMoney(targetUser.ID, amount)
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'argent de %s a été défini à %d.", targetUser.Username, amount))
					return
				}
			}
		}
	}
}
