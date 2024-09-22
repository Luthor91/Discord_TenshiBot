package utility_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

var commands = map[string]string{
	"affinity":           "Permet de voir l'affinité qu'on a avec le bot.",
	"burnaffinity":       "Brûler une affinité avec le bot.",
	"getaffinity":        "Voir son affinité avec le bot.",
	"setaffinity":        "Définir une affinité avec le bot.",
	"ban":                "Bannir un utilisateur avec une certaine raison.",
	"kick":               "Kick un utilisateur.",
	"delete":             "Supprimer un nombre spécifié de messages.",
	"timeout":            "Timeout quelqu'un pendant un moment.",
	"deafen":             "Déafen un utilisateur dans un canal vocal.",
	"kickvoice":          "Kick un utilisateur d'un canal vocal.",
	"mute":               "Muter un utilisateur dans un canal vocal.",
	"move":               "Déplacer un utilisateur dans un autre canal vocal.",
	"addgoodword":        "Ajouter un mot bon à la liste.",
	"addbadword":         "Ajouter un mot mauvais à la liste.",
	"deletegoodword":     "Supprimer un mot bon.",
	"deletebadword":      "Supprimer un mot mauvais.",
	"getgoodwords":       "Voir les mots bons.",
	"getbadwords":        "Voir les mots mauvais.",
	"createtextchannel":  "Créer un canal texte.",
	"createvoicechannel": "Créer un canal vocal.",
	"deletechannel":      "Supprimer un canal par son nom.",
	"burnexperience":     "Brûler des points d'expérience.",
	"getexperience":      "Voir son expérience.",
	"setxp":              "Définir son expérience.",
	"givexp":             "Donner de l'expérience à un utilisateur.",
	"leaderboard":        "Voir le leaderboard pour la monnaie.",
	"rank":               "Voir son classement pour la monnaie.",
	"burnmoney":          "Brûler de l'argent.",
	"getmoney":           "Voir combien d'argent on possède.",
	"daily":              "Récupérer de l'argent chaque jour.",
	"givemoney":          "Donner de l'argent à un utilisateur.",
	"setmoney":           "Définir une somme d'argent pour un utilisateur.",
	"calculate":          "Calculer une expression mathématique.",
	"ping":               "Voir la latence du bot.",
	"random":             "Générer un nombre aléatoire entre deux nombres.",
	"reminder":           "Créer un message timé.",
	"inventory":          "Voir son inventaire d'articles.",
	"use":                "Utiliser un item.",
	"giveitem":           "Donner un item à un utilisateur.",
	"getlogs":            "Récupérer les logs.",
	"getuserlogs":        "Récupérer les logs d'un utilisateur spécifique.",
	"championrotation":   "Voir la rotation des champions.",
	"summonerprofile":    "Voir le profil d'un invocateur.",
	"championinfo":       "Voir les informations d'un champion.",
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
