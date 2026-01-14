package banword_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

var wordService *services.WordService

func InitWordCommands(ws *services.WordService) {
	wordService = ws
}

// WordCommand gère les opérations sur les badwords
func WordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifie si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	// Définir le préfixe de commande
	command := fmt.Sprintf("%sbanword", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer et analyser les arguments de la commande
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Variables pour stocker l'état des options
	var addWord bool
	var deleteWord bool
	var listWords bool
	var verbose bool
	var specifiedWord string

	// Analyser les arguments
	for _, arg := range parsedArgs {
		switch arg.Arg {
			case "-a":
				addWord = true
			case "-d":
				deleteWord = true
			case "-l":
				listWords = true
			case "-v":
				verbose = true
			case "-h":
				showHelpMessage(s, m.ChannelID)
				return
		}
	}

	// Logique de commande
	if addWord && deleteWord {
		return
	}

	if addWord || deleteWord {
		specifiedWord = strings.TrimSpace(parsedArgs[len(parsedArgs)-1].Value)
	}

	if addWord {
		if specifiedWord == "" {
			return
		}
		handleBadWord(s, m, specifiedWord, verbose)
		return
	}

	if deleteWord {
		if specifiedWord == "" {
			return
		}
		handleDeleteWord(s, m, specifiedWord, verbose) // Appelle la fonction helper pour supprimer un mot
		return
	}

	if listWords {
		handleListBadWords(s, m) // Appelle la fonction pour lister les badwords
		return
	}

	if verbose {

		if addWord && specifiedWord == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à ajouter.")
		}
		if deleteWord && specifiedWord == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à supprimer.")
		}
		if addWord || deleteWord {
			s.ChannelMessageSend(m.ChannelID, "Vous ne pouvez pas ajouter et supprimer un mot en même temps.")
		}

		return
	}

	s.ChannelMessageSend(m.ChannelID, "Aucune commande reconnue. Usage : `?word [-g|-b|-d|-a|-l|-v|-h|-? <value>]`")
}
