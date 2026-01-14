package utility_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

var commands = map[string]string{
	"ban":                "Bannir un utilisateur avec une certaine raison.",
	"kick":               "Kick un utilisateur.",
	"delete":             "Supprimer un nombre spécifié de messages.",
	"timeout":            "Timeout quelqu'un pendant un moment.",
	"deafen":             "Mute un utilisateur dans un canal vocal.",
	"kickvoice":          "Kick un utilisateur d'un canal vocal.",
	"mute":               "Muter un utilisateur dans un canal vocal.",
	"move":               "Déplacer un utilisateur dans un autre canal vocal.",
	"banword":         	  "Ajouter ou retirer un banword.",
	"channel":      	  "Crééer, Supprime ou Archive un channel.",
	"statut":              "Voir le statut du serveur, user ou bot",
	"calculate":          "Calculer une expression mathématique.",
	"ping":               "Voir la latence du bot.",
	"random":             "Générer un nombre aléatoire entre deux nombres.",
	"reminder":           "Créer un message timé.",
	"logs":               "Récupérer les logs.",
}

// HelpCommand répond avec la liste des commandes ou des informations sur une commande spécifique.
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
