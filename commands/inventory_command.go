package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// InventoryCommand affiche l'inventaire de l'utilisateur
func InventoryCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sinventory", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer les items de l'utilisateur
	items, err := services.NewItemService().GetUserItems(m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de l'inventaire.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	// Vérifier si l'utilisateur a des items
	if len(items) == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Votre inventaire est vide.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	// Préparer le message de l'inventaire
	messageContent := "**Votre inventaire :**\n"
	for _, item := range items {
		messageContent += fmt.Sprintf("%s: %d\n", item.Name, item.Quantity)
	}

	_, err = s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil {
		log.Println("Erreur lors de l'envoi du message:", err)
	}
}
