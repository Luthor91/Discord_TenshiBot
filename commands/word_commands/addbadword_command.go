package word_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/database"
	"github.com/bwmarrin/discordgo"
)

// AddBadWordCommand ajoute un mot dans la liste des "badwords"
func AddBadWordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%saddbword", config.AppConfig.BotPrefix)
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
		if err := wordController.AddBadWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout du mot : "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Le mot a été ajouté avec succès.")
	}
}
