package word_commands

import (
	"fmt"

	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Création d'une instance de WordService
var wordService = services.NewWordService()

// handleGoodWord gère l'ajout d'un bon mot
func handleGoodWord(s *discordgo.Session, m *discordgo.MessageCreate, word string) {
	err := wordService.AddGoodWord(word)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de l'ajout du bon mot : %s", err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ajout du bon mot : %s", word))
}

// handleBadWord gère l'ajout d'un mauvais mot
func handleBadWord(s *discordgo.Session, m *discordgo.MessageCreate, word string) {
	err := wordService.AddBadWord(word)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de l'ajout du mauvais mot : %s", err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ajout du mauvais mot : %s", word))
}

// handleDeleteWord gère la suppression d'un mot
func handleDeleteWord(s *discordgo.Session, m *discordgo.MessageCreate, word string) {
	err := wordService.DeleteGoodWord(word)
	if err != nil {
		// Vérifiez si c'est un mauvais mot à supprimer si la suppression d'un bon mot échoue
		err = wordService.DeleteBadWord(word)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la suppression du mot : %s", err))
			return
		}
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Suppression du mot : %s", word))
}

// handleListGoodWords affiche les bons mots
func handleListGoodWords(s *discordgo.Session, m *discordgo.MessageCreate) {
	words, err := wordService.ListGoodWords()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des bons mots.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Voici la liste des bons mots : %v", words))
}

// handleListBadWords affiche les mauvais mots
func handleListBadWords(s *discordgo.Session, m *discordgo.MessageCreate) {
	words, err := wordService.ListBadWords()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des mauvais mots.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Voici la liste des mauvais mots : %v", words))
}
