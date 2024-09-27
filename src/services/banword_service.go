package services

import (
	"log"
	"strings"

	"github.com/Luthor91/DiscordBot/controllers"
	"github.com/Luthor91/DiscordBot/database"
	"github.com/bwmarrin/discordgo"
)

// BanWordService est un service pour gérer les mots interdits
type BanWordService struct {
	discordSession *discordgo.Session
	wordCtrl       *controllers.WordController
}

// NewBanWordService crée une nouvelle instance de BanWordService
func NewBanWordService(discordSession *discordgo.Session) *BanWordService {
	return &BanWordService{
		discordSession: discordSession,
		wordCtrl:       &controllers.WordController{DB: database.DB},
	}
}

// DeleteBanWordMessages supprime les messages contenant des mots interdits
func (service *BanWordService) DeleteBanWordMessages(m *discordgo.MessageCreate) {
	if m.Author.ID == service.discordSession.State.User.ID {
		return
	}

	// Récupérer la liste des mots interdits depuis le contrôleur
	badwords, err := service.wordCtrl.GetBadWords()
	if err != nil {
		log.Printf("Erreur lors de la récupération des mots interdits: %v", err)
		return
	}

	// Parcourir les mots interdits et vérifier s'ils sont présents dans le message
	for _, word := range badwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			// Supprimer le message
			err := service.discordSession.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.Printf("Erreur lors de la suppression du message contenant un mot interdit: %v", err)
			}
			break
		}
	}
}
