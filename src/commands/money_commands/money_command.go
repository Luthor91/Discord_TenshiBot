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

		if len(args) == 1 {
			handleShowUserBalance(s, m)
			return
		}

		if len(args) >= 2 && args[1] == "-h" {
			displayHelpMessage(s, m.ChannelID)
			return
		}

		isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)

		// Appeler handleTarget une fois ici si une cible est nécessaire
		var targetUser *discordgo.User
		if args[1] == "-n" || args[1] == "-g" || args[1] == "-s" || args[1] == "-m" {
			targetUser = handleTarget(s, m)
			if targetUser == nil {
				return
			}
		}

		switch args[1] {
		case "-r":
			if isMod {
				handleRemoveMoney(s, m, args)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		case "-d":
			handleDailyReward(s, m)
		case "-m":
			if isMod {
				handleShowMoney(s, m, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour voir le solde d'un autre utilisateur.")
			}
		case "-g":
			if isMod {
				handleGiveMoney(s, m, args, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		case "-s":
			if isMod {
				handleSetMoney(s, m, args, targetUser)
			} else {
				s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour effectuer cette commande.")
			}
		default:
			displayHelpMessage(s, m.ChannelID)
		}
	}
}
