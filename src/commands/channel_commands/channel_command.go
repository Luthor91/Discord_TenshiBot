package channel_commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// Fonction principale pour gérer la commande channel
func ChannelCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission d'exécuter cette commande.")
		return
	}

	command := fmt.Sprintf("%schannel", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%schan", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	// Si seul ?channel est utilisé, afficher les arguments possibles
	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Récupérer les arguments de la commande
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var channel *discordgo.Channel
	var duration time.Duration
	var isVoice bool
	var shouldLock, createChannelFlag, deleteChannelFlag bool
	var archiveMessagesCount int

	// Analyser les arguments extraits
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			var handleErr error
			channel, handleErr = discord.HandleChannel(s, m, arg.Value)
			if handleErr != nil {
				s.ChannelMessageSend(m.ChannelID, handleErr.Error())
				return
			}
		case "-t":
			duration = arg.Duration
		case "-v":
			isVoice = true
		case "-l":
			shouldLock = true
		case "-c":
			createChannelFlag = true
		case "-d":
			deleteChannelFlag = true
		case "-a":
			archiveMessagesCountValue, parseErr := strconv.Atoi(arg.Value)
			if parseErr == nil {
				archiveMessagesCount = archiveMessagesCountValue
			}
		default:
			showHelpMessage(s, m.ChannelID)
			return
		}
	}

	// Créer un salon si l'option -c est présente
	if createChannelFlag {
		// Vérifiez que le canal est bien mentionné
		if channel == nil {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nom de salon valide.")
			return
		}

		err := discord.CreateChannel(s, m.GuildID, channel.ID, channel.Name, isVoice, duration)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la création du salon : "+err.Error())
		}
	}

	// Supprimer un salon si l'option -d est présente
	if deleteChannelFlag {
		err := discord.DeleteChannel(s, m.GuildID, channel.ID, channel.Name)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du salon : "+err.Error())
		}
	}

	// Verrouiller ou déverrouiller un salon si l'option -l est présente
	if shouldLock {
		err := discord.HandleChannelLock(s, m, channel.ID, duration)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la gestion du verrouillage du salon : "+err.Error())
		}
	}

	// Archive messages si l'option -a est présente
	if archiveMessagesCount > 0 {
		err := archiveMessages(s, m, archiveMessagesCount)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'archivage des messages : "+err.Error())
		}
	}
}
