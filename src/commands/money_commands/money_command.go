package money_commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/Luthor91/DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
)

// MoneyCommand gère les commandes liées à l'argent
func MoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Indicateur pour savoir si les réponses doivent être affichées
	showResponses := false
	for _, arg := range parsedArgs {
		if arg.Arg == "-v" {
			showResponses = true
			break // Sortie une fois qu'on a trouvé -v
		}
	}

	// Si aucune option n'est spécifiée, retourner le solde de l'utilisateur
	if len(parsedArgs) == 0 {
		money, err := services.NewUserService().GetMoney(m.Author.ID)
		if err != nil {
			utils.SendResponse(s, m.ChannelID, "Erreur lors de la récupération du solde.", showResponses)
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre money est de %d.", m.Author.Username, money))
		return
	}

	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	var targetUser *discordgo.User

	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			targetUser = discord.HandleTarget(s, m, arg.Value)
			if targetUser == nil {
				return
			}
		case "-r":
			if isMod || targetUser == m.Author {
				if arg.Value != "" {
					amount, err := strconv.Atoi(arg.Value)
					if err != nil || amount <= 0 {
						utils.SendResponse(s, m.ChannelID, "Veuillez entrer un montant valide pour retirer.", showResponses)
						return
					}

					services.NewUserService().AddMoney(m.Author.ID, -amount) // Retirer de l'argent
					utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous avez retiré %d unités.", amount), showResponses)
				} else {
					utils.SendResponse(s, m.ChannelID, "Veuillez spécifier un montant à retirer.", showResponses)
				}
			} else {
				utils.SendResponse(s, m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.", showResponses)
			}
		case "-d":
			canReceive, timeLeft, err := services.NewUserService().CanReceiveDailyReward(m.Author.ID)
			if !canReceive || err != nil {
				utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous devez attendre encore %v avant de réclamer la prochaine récompense quotidienne.", timeLeft.Round(time.Minute)), showResponses)
				return
			}
			randomAmount := rand.Intn(91) + 10
			services.NewUserService().UpdateDailyMoney(m.Author.ID, randomAmount)
			utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous avez reçu %d unités aujourd'hui !", randomAmount), showResponses)
		case "-m":
			if isMod || targetUser == m.Author {
				money, err := services.NewUserService().GetMoney(targetUser.ID)
				if err != nil {
					utils.SendResponse(s, m.ChannelID, "Erreur lors de la récupération du solde.", showResponses)
					return
				}
				utils.SendResponse(s, m.ChannelID, fmt.Sprintf("%s a %d unités de monnaie.", targetUser.Username, money), showResponses)
			} else {
				utils.SendResponse(s, m.ChannelID, "Vous n'avez pas la permission pour voir le solde d'un autre utilisateur.", showResponses)
			}
		case "-g":
			if targetUser != nil {
				amount, err := strconv.Atoi(arg.Value)
				if err != nil || amount <= 0 {
					utils.SendResponse(s, m.ChannelID, "Veuillez entrer un montant valide pour donner.", showResponses)
					return
				}

				services.NewUserService().GiveMoney(m.Author.ID, targetUser.ID, amount)
				utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous avez donné %d unités à %s.", amount, targetUser.Username), showResponses)
			} else {
				utils.SendResponse(s, m.ChannelID, "Veuillez spécifier un utilisateur à qui donner de l'argent.", showResponses)
			}
		case "-s":
			if isMod {
				amount, err := strconv.Atoi(arg.Value)
				if err != nil || amount < 0 {
					utils.SendResponse(s, m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.", showResponses)
					return
				}

				services.NewUserService().SetMoney(targetUser.ID, amount)
				utils.SendResponse(s, m.ChannelID, fmt.Sprintf("L'argent de %s a été défini à %d.", targetUser.Username, amount), showResponses)
			} else {
				utils.SendResponse(s, m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.", showResponses)
			}
		case "-a":
			if isMod {
				amount, err := strconv.Atoi(arg.Value)
				if err != nil || amount <= 0 {
					utils.SendResponse(s, m.ChannelID, "Veuillez entrer un montant valide à ajouter.", showResponses)
					return
				}

				if targetUser != nil {
					services.NewUserService().AddMoney(targetUser.ID, amount)
					utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous avez ajouté %d unités à %s.", amount, targetUser.Username), showResponses)
				} else {
					services.NewUserService().AddMoney(m.Author.ID, amount)
					utils.SendResponse(s, m.ChannelID, fmt.Sprintf("Vous avez ajouté %d unités à votre solde.", amount), showResponses)
				}
			} else {
				utils.SendResponse(s, m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.", showResponses)
			}
		case "-h":
			displayHelpMessage(s, m.ChannelID)
		default:
			utils.SendResponse(s, m.ChannelID, "Commande inconnue. Utilisez -h pour voir l'aide.", showResponses)
		}
	}
}

// Fonction pour afficher un message d'aide
func displayHelpMessage(s *discordgo.Session, channelID string) {
	message := "Utilisation de la commande money:\n" +
		"-n : Utiliser le nom d'utilisateur\n" +
		"-r [montant] : Retirer de l'argent\n" +
		"-d : Récompense quotidienne\n" +
		"-m : Afficher votre solde\n" +
		"-g [montant] : Donner de l'argent à un utilisateur\n" +
		"-s [montant] : Définir l'argent d'un utilisateur (admin uniquement)\n" +
		"-a [montant] : Ajouter de l'argent à un utilisateur (admin uniquement)\n" +
		"-h : Afficher ce message d'aide\n" +
		"-v : Activer l'affichage des réponses du bot pour les commandes"
	s.ChannelMessageSend(channelID, message)
}
