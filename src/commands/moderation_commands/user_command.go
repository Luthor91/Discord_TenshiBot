package moderation_commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// ModerateUserCommand gère les commandes de modération des utilisateurs
func ModerateUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour exécuter cette commande.")
		return
	}

	// Vérifier si la commande commence par le bon préfixe
	command := fmt.Sprintf("%suser", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%susr", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	// Extraction des arguments
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Si aucun argument n'est renseigné, afficher un message explicatif
	if len(parsedArgs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Arguments possibles :\n"+
			"-n [user]: mentionner l'utilisateur\n"+
			"-b [reason]: bannir l'utilisateur avec une raison\n"+
			"-w [reason]: avertir l'utilisateur avec une raison\n"+
			"-k [reason]: expulser l'utilisateur avec une raison\n"+
			"-m [duration]: mettre l'utilisateur en sourdine pour une durée\n"+
			"-d [duration]: rendre l'utilisateur sourd pour une durée\n"+
			"-to [duration]: mettre l'utilisateur en timeout pour une durée\n"+
			"-mv [channel]: déplacer l'utilisateur dans un canal spécifique\n"+
			"-t [duration]: durée de l'action (pour mute, deafen, timeout)\n"+
			"-r: réinitialiser tous les statuts de l'utilisateur\n"+
			"-rw: réinitialiser les avertissements de l'utilisateur")
		return
	}

	// Récupérer les valeurs des parsedArgs
	var userID string
	var action string
	var reason string
	var actionTime time.Duration
	var targetChannel string

	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			userID = discord.HandleTarget(s, m, arg.Value).ID
		case "-b":
			action = "ban"
			reason = arg.Value
		case "-w":
			action = "warn"
			reason = arg.Value
		case "-k":
			action = "kick"
			reason = arg.Value
		case "-m":
			action = "mute"
			reason = arg.Value
		case "-d":
			action = "deafen"
			reason = arg.Value
		case "-to":
			action = "timeout"
			reason = arg.Value
		case "-mv":
			action = "move"
			targetChannel = arg.Value
		case "-t":
			actionTime = arg.Duration
		case "-r":
			if userID != "" {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("(En cours) Statut de l'utilisateur %s réinitialisé.", userID))
				return
			}
		case "-rw":
			if userID != "" {
				warnService := services.NewWarnService(s, m.GuildID)
				err := warnService.ResetWarns(userID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la réinitialisation des avertissements de l'utilisateur %s: %v", userID, err))
					return
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissements de l'utilisateur %s réinitialisés.", userID))
				return
			}
		default:
			// Ignorer les arguments non reconnus
		}
	}

	// Effectuer l'action spécifiée
	switch action {
	case "ban":
		discord.BanUser(s, m, userID, reason)
	case "warn":
		warnService := services.NewWarnService(s, m.GuildID)
		err := warnService.AddWarn(userID, reason, m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de l'envoi de l'avertissement à l'utilisateur %s : %v", userID, err))
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Avertissement envoyé à l'utilisateur %s : %s", userID, reason))
	case "kick":
		discord.KickUser(s, m, userID, reason)
	case "mute":
		discord.MuteUser(s, m, userID, actionTime, reason)
	case "deafen":
		discord.DeafenUser(s, m, userID, actionTime, reason)
	case "timeout":
		discord.TimeoutUser(s, m.GuildID, userID, actionTime)
	case "move":
		discord.MoveUser(s, m, userID, targetChannel)
	default:
		s.ChannelMessageSend(m.ChannelID, "Action inconnue.")
	}
}
