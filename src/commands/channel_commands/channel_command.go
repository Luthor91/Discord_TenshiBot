package channel_commands

import (
	"fmt"
	"strings"

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
	channelName, duration, isVoice, shouldLock, createChannelFlag, deleteChannelFlag, archiveMessagesCount, err := parseChannelArgs(m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Créer un salon si l'option -c est présente
	if createChannelFlag {
		err := discord.CreateChannel(s, m.GuildID, m.ChannelID, channelName, isVoice, duration)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la création du salon : "+err.Error())
		}
	}

	// Supprimer un salon si l'option -d est présente
	if deleteChannelFlag {
		err := discord.DeleteChannel(s, m.GuildID, m.ChannelID, channelName)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du salon : "+err.Error())
		}
	}

	// Verrouiller ou déverrouiller un salon si l'option -l est présente
	if shouldLock {
		err := discord.HandleChannelLock(s, m, m.ChannelID, duration)
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
