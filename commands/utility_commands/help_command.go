package utility_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

var commands = map[string]string{
	"affinity":    "Permet de voir l'affinité qu'on a avec le bot.",
	"ban":         "Bannir un utilisateur avec une certaine raison.",
	"kick":        "Kick un utilisateur.",
	"calculate":   "Calculer une expression mathématique.",
	"daily":       "Récupérer de l'argent chaque jour.",
	"delete":      "Supprimer un nombre spécifié de messages.",
	"xp":          "Voir son XP.",
	"leaderboard": "Voir le leaderboard pour la monnaie.",
	"money":       "Voir combien d'argent on possède.",
	"ping":        "Voir la latence du bot.",
	"random":      "Générer un nombre aléatoire entre deux nombres.",
	"rank":        "Voir son classement pour la monnaie.",
	"reminder":    "Créer un message timé.",
	"timeout":     "Timeout quelqu'un pendant un moment.",
}

// HelpCommand répond "Hello World😃" lorsque l'utilisateur tape "!help"
func HelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%shelp", config.AppConfig.BotPrefix)

	if m.Content == command {
		args := strings.Fields(m.Content)
		prefix := config.AppConfig.BotPrefix

		// Si aucun argument supplémentaire, lister les commandes
		if len(args) == 1 {
			commandList := "Liste des commandes :\n"
			for cmd := range commands {
				commandList += fmt.Sprintf("`%s%s`\n", prefix, cmd)
			}
			s.ChannelMessageSend(m.ChannelID, commandList)
			return
		}

		// Si un argument est passé, donner plus d'infos sur la commande
		command := strings.ToLower(args[1])
		if description, exists := commands[command]; exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s%s` : %s", prefix, command, description))
		} else {
			s.ChannelMessageSend(m.ChannelID, "Commande non reconnue.")
		}
	}
}
