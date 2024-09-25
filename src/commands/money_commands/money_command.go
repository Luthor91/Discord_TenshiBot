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
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)

		// Si seulement ?money est renseigné, afficher le solde de l'utilisateur
		if len(args) == 1 {
			handleShowUserBalance(s, m)
			return
		}

		if len(args) >= 2 && args[1] == "-h" {
			displayHelpMessage(s, m.ChannelID)
			return
		}

		isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)

		// Gérer les commandes spéciales avec ou sans cible
		var targetUser *discordgo.User
		if len(args) >= 2 && args[1] == "-n" {
			// Si -n est spécifié, récupérer la cible
			targetUser = handleTarget(s, m)
			if targetUser == nil {
				return
			}
		} else {
			// Sinon, prendre l'auteur comme cible par défaut
			targetUser = m.Author
		}

		// Vérifier le type de commande fourni après -n ou pour l'auteur
		switch args[1] {
		case "-r":
			// Commande pour retirer de l'argent (modérateur uniquement)
			if isMod {
				handleRemoveMoney(s, m, args)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		case "-d":
			// Réclamer la récompense quotidienne
			handleDailyReward(s, m)
		case "-m":
			// Afficher l'argent d'un utilisateur (modérateur uniquement)
			if isMod {
				handleShowMoney(s, m, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour voir le solde d'un autre utilisateur.")
			}
		case "-g":
			// Donner de l'argent à un utilisateur (modérateur uniquement)
			if isMod {
				handleGiveMoney(s, m, args, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		case "-s":
			// Définir l'argent de l'utilisateur, soit l'auteur soit la cible (-n)
			if len(args) >= 3 {
				handleSetMoney(s, m, args, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous devez spécifier un montant pour définir l'argent.")
			}
		default:
			displayHelpMessage(s, m.ChannelID)
		}
	}
}
