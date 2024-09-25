package services

import (
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/bwmarrin/discordgo"
)

// MessageService est un service pour gérer les messages
type MessageService struct {
	userService     *UserService
	affinityService *AffinityService
	logService      *LogService
}

type KeywordResponse struct {
	Keyword  string
	Response string
	Reaction string
}

var keywordResponsesWithMention = []KeywordResponse{
	{Keyword: "help", Response: "Je suis ici pour t'aider !", Reaction: "✅"},
	// Ajoutez d'autres mots-clés qui nécessitent une mention
}

var keywordResponsesWithoutMention = []KeywordResponse{
	{Keyword: "bonjour", Response: "Salut à toi !", Reaction: "👋"},
	{Keyword: "merci", Response: "De rien !", Reaction: "😊"},
	{Keyword: "aide", Response: "Voici comment je peux t'aider...", Reaction: "❓"},
	{Keyword: "gg", Response: "", Reaction: "👏"},                   // Juste une réaction
	{Keyword: "lol", Response: "Haha, très drôle !", Reaction: ""}, // Juste une réponse
}

// NewMessageService crée une nouvelle instance de MessageService
func NewMessageService(userService *UserService, affinityService *AffinityService, logService *LogService) *MessageService {
	return &MessageService{
		userService:     userService,
		affinityService: affinityService,
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

	user, err := controllers.NewUserController().GetUserByDiscordID(message.Author.ID)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur : %v", err)
		return
	}

	if err := service.userService.AddExperience(user.UserDiscordID, 1); err != nil {
		log.Printf("Erreur lors de l'ajout de l'expérience : %v", err)
		return
	}

	service.affinityService.AdjustAffinity(user.UserDiscordID, message)

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
		if strings.Contains(content, keywordResponse.Keyword) {
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

			// Sortir de la boucle une fois que le mot-clé est trouvé
			return
		}
	}
}
