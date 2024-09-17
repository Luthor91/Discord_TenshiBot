package features

import (
	"github.com/bwmarrin/discordgo"
)

// newServerMessage est appelé lorsqu'un nouveau message est reçu dans un serveur.
func NewServerMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Éviter que le bot réponde à ses propres messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Ajouter de la monnaie et de l'expérience à l'utilisateur
	AddUserIfNotExists(message.Author.ID, message.Author.Username)
	AddExperience(message.Author.ID, 1)
	AdjustAffinity(discord, message)

	// Logger le message
	LogMessage(discord, message)
}

// newPrivateMessage est appelé lorsqu'un nouveau message est reçu en message privé (DM).
func NewPrivateMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// Éviter de répondre aux propres messages du bot ou aux messages publics
	if message.Author.ID == discord.State.User.ID || message.GuildID != "" {
		return
	}

	// Répondre au message privé
	discord.ChannelMessageSend(message.ChannelID, "Merci pour votre message privé !")
}
