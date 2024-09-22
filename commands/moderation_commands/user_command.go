package moderation_commands

import (
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

func ModerateUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur (adapter cette vérification selon votre logique)
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour exécuter cette commande.")
		return
	}

	// Vérifier si la commande commence par "?user"
	if !strings.HasPrefix(m.Content, "?user") {
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

	// Retrieving values from parsedArgs
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
		banUser(s, m, userID, reason)
	case "warn":
		warnUser(s, m, userID, reason)
	case "kick":
		kickUser(s, m, userID, reason)
	case "mute":
		muteUser(s, m, userID, actionTime, reason)
	case "deafen":
		deafenUser(s, m, userID, actionTime, reason)
	case "timeout":
		timeoutUser(s, m, userID, actionTime, reason)
	case "move":
		moveUser(s, m, userID, targetChannel)
	default:
		s.ChannelMessageSend(m.ChannelID, "Action inconnue.")
	}
}
