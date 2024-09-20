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

// GetBadWordsCommand liste tous les mauvais mots
func GetBadWordsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sgetbword", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	wordController := &controllers.WordController{
		DB: database.DB,
	}

	badWords, err := wordController.GetBadWords()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des mauvais mots : "+err.Error())
		return
	}

	if len(badWords) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Aucun mauvais mot trouvé.")
		return
	}

	response := "Liste des mauvais mots :\n" + strings.Join(badWords, "\n")
	s.ChannelMessageSend(m.ChannelID, response)
}
