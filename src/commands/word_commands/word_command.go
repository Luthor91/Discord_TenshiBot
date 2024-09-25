package word_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// WordCommand gère les opérations sur les mots (goodwords et badwords)
func WordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifie si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires.")
		return
	}

	// Définir le préfixe de commande
	command := fmt.Sprintf("%sword", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer et analyser les arguments de la commande
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Vérifier qu'il y a des arguments
	if len(parsedArgs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Usage : `?word [-g|-b|-d|-a|-l|-v|-h|-? <value>]`")
		return
	}

	// Variables pour stocker l'état des options
	var goodWord bool
	var badWord bool
	var addWord bool
	var deleteWord bool
	var listWords bool
	var verbose bool
	var specifiedWord string

	// Analyser les arguments
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-g":
			goodWord = true
		case "-b":
			badWord = true
		case "-a":
			addWord = true
		case "-d":
			deleteWord = true
		case "-l":
			listWords = true
		case "-v":
			verbose = true
		case "-h":
			s.ChannelMessageSend(m.ChannelID, "Usage : `?word [-g|-b|-d|-a|-l|-v|-h|-? <value>]`")
			return
		case "-?":
			if len(arg.Value) > 0 {
				specifiedWord = arg.Value
			}
		default:
			s.ChannelMessageSend(m.ChannelID, "Argument non reconnu. Usage : `?word [-g|-b|-d|-a|-l|-v|-h|-? <value>]`")
			return
		}
	}

	// Logique de commande
	if addWord && deleteWord {
		s.ChannelMessageSend(m.ChannelID, "Vous ne pouvez pas ajouter et supprimer un mot en même temps.")
		return
	}

	if addWord {
		if specifiedWord == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à ajouter avec `-? <value>`.")
			return
		}
		if goodWord {
			handleGoodWord(s, m, specifiedWord) // Appelle la fonction helper pour gérer les goodwords
		} else if badWord {
			handleBadWord(s, m, specifiedWord) // Appelle la fonction helper pour gérer les badwords
		} else {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier si le mot est un bon mot (-g) ou un mauvais mot (-b).")
		}
		return
	}

	if deleteWord {
		if specifiedWord == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à supprimer avec `-? <value>`.")
			return
		}
		handleDeleteWord(s, m, specifiedWord) // Appelle la fonction helper pour supprimer un mot
		return
	}

	if listWords {
		if goodWord {
			handleListGoodWords(s, m) // Appelle la fonction pour lister les goodwords
		}
		if badWord {
			handleListBadWords(s, m) // Appelle la fonction pour lister les badwords
		}
		return
	}

	if verbose {
		// Ajoutez ici la logique pour le mode verbose si nécessaire
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Aucune commande reconnue. Usage : `?word [-g|-b|-d|-a|-l|-v|-h|-? <value>]`")
}
