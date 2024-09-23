package services

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// InteractionService gère les interactions spécifiques avec le bot
type InteractionService struct {
	BotID string
}

// NewInteractionService crée une nouvelle instance du service d'interaction
func NewInteractionService(botID string) *InteractionService {
	return &InteractionService{BotID: botID}
}

// HandleMessage gère les messages envoyés sur le serveur pour des interactions spécifiques
func (s *InteractionService) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Ignore les messages envoyés par le bot
	if message.Author.ID == s.BotID {
		return
	}

	// Vérifie si le bot a été mentionné dans le message
	if strings.Contains(message.Content, "<@"+s.BotID+">") {
		s.handleBotMention(session, message)
		return
	}

	// Vérifie des mots-clés spécifiques
	keywords := []string{"hello", "help", "info"}
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(message.Content), keyword) {
			s.handleKeyword(session, message, keyword)
			return
		}
	}
}

// handleBotMention est appelé lorsque le bot est mentionné dans un message
func (s *InteractionService) handleBotMention(session *discordgo.Session, message *discordgo.MessageCreate) {
	response := "Salut! Tu m'as mentionné, comment puis-je t'aider ?"
	session.ChannelMessageSend(message.ChannelID, response)
}

// handleKeyword gère des mots-clés spécifiques dans les messages
func (s *InteractionService) handleKeyword(session *discordgo.Session, message *discordgo.MessageCreate, keyword string) {
	switch keyword {
	case "hello":
		session.ChannelMessageSend(message.ChannelID, "Bonjour!")
	case "help":
		session.ChannelMessageSend(message.ChannelID, "Voici la liste des commandes disponibles...")
	case "info":
		session.ChannelMessageSend(message.ChannelID, "Je suis un bot Discord, prêt à t'aider!")
	}
}
