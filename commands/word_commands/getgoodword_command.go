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

// GetGoodWordsCommand liste tous les bons mots
func GetGoodWordsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%sgetgword", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	wordController := &controllers.WordController{
		DB: database.DB,
	}

	goodWords, err := wordController.GetGoodWords()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des bons mots : "+err.Error())
		return
	}

	if len(goodWords) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Aucun bon mot trouvé.")
		return
	}

	response := "Liste des bons mots :\n" + strings.Join(goodWords, "\n")
	s.ChannelMessageSend(m.ChannelID, response)
}
