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

// DeleteBadWordCommand supprime un mot de la liste des "badwords"
func DeleteBadWordCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sdelbword", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, command) {
		word := strings.TrimSpace(strings.TrimPrefix(m.Content, command))
		if word == "" {
			s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un mot à supprimer.")
			return
		}

		wordController := &controllers.WordController{
			DB: database.DB,
		}

		if err := wordController.DeleteBadWord(word); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la suppression du mot : "+err.Error())
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Le mot a été supprimé avec succès.")
	}
}
