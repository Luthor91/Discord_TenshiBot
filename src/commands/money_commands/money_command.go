package money_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

func MoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := fmt.Sprintf("%smoney", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)

		// Vérification d'arguments vides
		if len(args) < 2 {
			displayHelpMessage(s, m.ChannelID)
			return
		}

		switch args[1] {
		case "-n":
			handleMention(s, m)
		case "-r":
			handleRemoveMoney(s, m, args)
		case "-d":
			handleDailyReward(s, m)
		case "-m":
			handleShowMoney(s, m)
		case "-g":
			handleGiveMoney(s, m, args)
		case "-s":
			handleSetMoney(s, m, args)
		default:
			displayHelpMessage(s, m.ChannelID)
		}
	}
}
