package commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/bwmarrin/discordgo"
)

// AddGoodWordCommand ajoute un mot dans la liste des "goodwords"
func AddGoodWordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%saddgoodword", config.AppConfig.BotPrefix)
	// Vérifie si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Récupère le mot à ajouter
		word := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if word == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à ajouter.")
			return
		}

		// Crée une instance de WordController
		wordController := &controllers.WordController{
			DB: database.DB, // Assurez-vous que vous avez configuré la connexion à la base de données dans la config
		}

		// Ajoute le mot dans la liste des "badwords"
		if err := wordController.AddGoodWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout du mot : "+err.Error())
			return
		}

		// Confirme l'ajout du mot
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Le mot \"%s\" a été ajouté à la liste des mots positifs.", word))
	}
}
