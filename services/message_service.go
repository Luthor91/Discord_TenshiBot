package services

import (
	"log"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/bwmarrin/discordgo"
)

// MessageService est un service pour gérer les messages
type MessageService struct {
	userService     *UserService
	affinityService *AffinityService
	logService      *LogService
}

// NewMessageService crée une nouvelle instance de MessageService
func NewMessageService(userService *UserService, affinityService *AffinityService, logService *LogService) *MessageService {
	return &MessageService{
		userService:     userService,
		affinityService: affinityService,
		logService:      logService,
	}
}

// NewServerMessage est appelé lorsqu'un nouveau message est reçu dans un serveur.
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

	if err := service.userService.AddExperience(user, 1); err != nil {
		log.Printf("Erreur lors de l'ajout de l'expérience : %v", err)
		return
	}

	service.affinityService.AdjustAffinity(message)

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
