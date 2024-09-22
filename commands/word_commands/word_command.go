package word_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
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

	// Récupérer les arguments de la commande
	args := strings.Fields(strings.TrimPrefix(m.Content, command))
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage : `?word [-g|-b|-d|-a|-l] [word]`")
		return
	}

	action := args[0]
	word := strings.Join(args[1:], " ")

	wordController := controllers.NewWordController()

	switch action {
	case "-g": // Ajouter un "goodword"
		if err := wordController.AddGoodWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout du bon mot : "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été ajouté aux goodwords.", word))

	case "-b": // Ajouter un "badword"
		if err := wordController.AddBadWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout du mauvais mot : "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été ajouté aux badwords.", word))

	case "-d": // Supprimer un mot
		if strings.HasPrefix(word, "good") {
			if err := wordController.DeleteGoodWord(word); err != nil {
				s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du goodword : "+err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été supprimé des goodwords.", word))
		} else if strings.HasPrefix(word, "bad") {
			if err := wordController.DeleteBadWord(word); err != nil {
				s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du badword : "+err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été supprimé des badwords.", word))
		} else {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier si le mot est un goodword ou un badword.")
		}

	case "-a": // Ajouter un mot (par défaut goodword)
		if err := wordController.AddGoodWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout du mot : "+err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été ajouté avec succès.", word))

	case "-l": // Lister les mots
		if word == "good" {
			goodWords, err := wordController.GetGoodWords()
			if err != nil || len(goodWords) == 0 {
				s.ChannelMessageSend(m.ChannelID, "Aucun goodword trouvé.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Liste des goodwords :\n"+strings.Join(goodWords, "\n"))
		} else if word == "bad" {
			badWords, err := wordController.GetBadWords()
			if err != nil || len(badWords) == 0 {
				s.ChannelMessageSend(m.ChannelID, "Aucun badword trouvé.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Liste des badwords :\n"+strings.Join(badWords, "\n"))
		} else {
			s.ChannelMessageSend(m.ChannelID, "Spécifiez 'good' ou 'bad' pour lister les mots correspondants.")
		}

	default:
		s.ChannelMessageSend(m.ChannelID, "Argument invalide. Usage : `?word [-g|-b|-d|-a|-l] [word]`")
	}
}
