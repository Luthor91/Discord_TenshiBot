package services

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MessageService est un service pour gérer les messages
type MessageService struct {
	userService     *UserService
	logService      *LogService
}

type KeywordResponse struct {
	Keyword  string
	Response string
	Reaction string
}

var keywordResponsesWithMention = []KeywordResponse{
	{Keyword: "gentil bot", Reaction: "😳"},
	{Keyword: "good bot", Response: "", Reaction: "😳"},
	{Keyword: "salut", Response: "Salut !", Reaction: "👋"},
	{Keyword: "bonjour", Response: "Bonjour !", Reaction: "👋"},
	{Keyword: "hey", Response: "Hey !", Reaction: "👋"},
}


var keywordResponsesWithoutMention = []KeywordResponse{
	{Keyword: "ntm", Response: "", Reaction: "😠"},
	{Keyword: "gg", Response: "", Reaction: "👏"},
	{Keyword: "belle bite", Response: "", Reaction: "👑"},
	{Keyword: "luthor", Response: "Oui c'est moi"},
	
}

// NewMessageService crée une nouvelle instance de MessageService
func NewMessageService(userService *UserService, logService *LogService) *MessageService {
	return &MessageService{
		userService:     userService,
		logService:      logService,
	}
}

func (service *MessageService) NewServerMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Éviter que le bot réponde à ses propres messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Ajouter de la monnaie et de l'expérience à l'utilisateur
	if err := service.userService.AddUserIfNotExists(message.Author.ID, message.Author.Username); err != nil {
		log.Printf("Erreur lors de l'ajout de l'utilisateur : %v", err)
		return
	}

	// Vérifier si le bot est mentionné
	if len(message.Mentions) > 0 {
		service.handleKeywordResponses(discord, message, true)
	} else {
		service.handleKeywordResponses(discord, message, false)
	}

	// Enregistrer le message dans les logs
	if err := service.logService.LogMessage(discord, message); err != nil {
		log.Printf("Erreur lors de l'enregistrement du message : %v", err)
		return
	}
}

// NewPrivateMessage est appelé lorsqu'un nouveau message est reçu en message privé (DM).
func (service *MessageService) NewPrivateMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Éviter de répondre aux propres messages du bot ou aux messages publics
	if message.Author.ID == discord.State.User.ID || message.GuildID != "" {
		return
	}

	// Répondre au message privé
	if _, err := discord.ChannelMessageSend(message.ChannelID, "Merci pour votre message privé !"); err != nil {
		log.Printf("Erreur lors de l'envoi de la réponse au message privé : %v", err)
	}
}

func (service *MessageService) handleKeywordResponses(discord *discordgo.Session, message *discordgo.MessageCreate, mentioned bool) {
	content := strings.ToLower(message.Content) // Normaliser le message en minuscule

	var keywordResponses []KeywordResponse
	if mentioned {
		keywordResponses = keywordResponsesWithMention
	} else {
		keywordResponses = keywordResponsesWithoutMention
	}

	for _, keywordResponse := range keywordResponses {
		pattern := `\b` + regexp.QuoteMeta(keywordResponse.Keyword) + `\b`
		matched, err := regexp.MatchString(pattern, content)
		if err != nil {
			log.Printf("Erreur regex : %v", err)
			continue
		}

		if matched {
			// Envoyer la réponse si elle n'est pas vide
			if keywordResponse.Response != "" {
				_, err := discord.ChannelMessageSend(message.ChannelID, keywordResponse.Response)
				if err != nil {
					log.Printf("Erreur lors de l'envoi de la réponse : %v", err)
				}
			}

			// Ajouter la réaction si elle n'est pas vide
			if keywordResponse.Reaction != "" {
				err := discord.MessageReactionAdd(message.ChannelID, message.ID, keywordResponse.Reaction)
				if err != nil {
					log.Printf("Erreur lors de l'ajout de la réaction : %v", err)
				}
			}

			return
		}
	}

}
