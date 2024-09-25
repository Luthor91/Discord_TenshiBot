package money_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// MoneyCommand gère les commandes liées à l'argent
func MoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Vérifie que le message ne provient pas du bot lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Préfixe de commande pour l'argent
	command := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Si aucun argument n'est fourni, afficher le solde de l'utilisateur
	if len(parsedArgs) == 0 {
		handleShowUserBalance(s, m)
		return
	}

	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)

	// Gérer les commandes spéciales avec ou sans cible
	var targetUser *discordgo.User
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			targetUser = discord.HandleTarget(s, m, arg.Value) // Passez la valeur de l'argument
			if targetUser == nil {
				return
			}
		case "-r":
			// Retirer de l'argent
			if isMod || targetUser == m.Author {
				// Vérifiez qu'il y a un montant spécifié après -r

				if arg.Value != "" {
					handleRemoveMoney(s, m, arg.Value) // Passer le montant
				} else {
					s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un montant à retirer.")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		case "-d":
			handleDailyReward(s, m)
		case "-m":
			if isMod || targetUser == m.Author {
				handleShowMoney(s, m, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour voir le solde d'un autre utilisateur.")
			}
		case "-g":
			if targetUser != nil {
				handleGiveMoney(s, m, arg.Value, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un utilisateur à qui donner de l'argent.")
			}
		case "-s":
			if isMod {
				handleSetMoney(s, m, arg.Value, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		default:
			displayHelpMessage(s, m.ChannelID)
		}
	}
}
