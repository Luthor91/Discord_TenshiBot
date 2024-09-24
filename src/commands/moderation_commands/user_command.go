package moderation_commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

func ModerateUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour exécuter cette commande.")
		return
	}

	// Vérifier si la commande
	command := fmt.Sprintf("%suser", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%susr", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}
	// Parsing command
	args := strings.Fields(m.Content)

	// Parse command arguments
	parsedArgs, err := parseArgs(args)
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

	// Récupérer les valeurs de parsedArgs
	userID := parsedArgs["-n"]
	action := ""
	reason := ""
	actionTime := time.Duration(0)
	targetChannel := ""

	// Gérer les différentes actions (ban, warn, kick, etc.)
	if val, exists := parsedArgs["-b"]; exists {
		action = "ban"
		reason = val
	}
	if val, exists := parsedArgs["-w"]; exists {
		action = "warn"
		reason = val
	}
	if val, exists := parsedArgs["-k"]; exists {
		action = "kick"
		reason = val
	}
	if val, exists := parsedArgs["-m"]; exists {
		action = "mute"
		reason = val
	}
	if val, exists := parsedArgs["-d"]; exists {
		action = "deafen"
		reason = val
	}
	if val, exists := parsedArgs["-to"]; exists {
		action = "timeout"
		reason = val
	}
	if val, exists := parsedArgs["-mv"]; exists {
		action = "move"
		targetChannel = val
	}
	if val, exists := parsedArgs["-t"]; exists {
		duration, err := utils.ParseDuration(val)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Durée incorrecte, veuillez vérifier.")
			return
		}
		actionTime = duration
	}
	if _, exists := parsedArgs["-r"]; exists {
		resetAllUserStatus(s, m, userID)
		return
	}
	if _, exists := parsedArgs["-rw"]; exists {
		resetUserWarnings(s, m, userID)
		return
	}

	// Effectuer l'action spécifiée
	switch action {
	case "ban":
		discord.BanUser(s, m, userID, reason)
	case "warn":
		warnUser(s, m, userID, reason)
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
