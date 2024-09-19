package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// UseItemCommand permet d'utiliser un item sur un autre utilisateur
func UseItemCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%suse", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments de la commande
	args := strings.Fields(m.Content[len(command):])
	if len(args) < 3 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: !use <item_name> @user <quantity>")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	itemName := args[0]
	targetMention := args[1]
	quantity, err := strconv.Atoi(args[2])
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "La quantité doit être un nombre.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	// Extraire l'ID de l'utilisateur cible depuis la mention
	targetID := strings.TrimPrefix(targetMention, "<@!")
	targetID = strings.TrimSuffix(targetID, ">")

	// Vérifier si l'utilisateur émetteur possède l'item
	hasItem, err := services.HasItem(m.Author.ID, itemName, quantity)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "Erreur lors de la vérification de l'item.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}
	if !hasItem {
		_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de cet item.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	// Utiliser l'item
	err = services.UseItem(m.Author.ID, targetID, itemName, quantity)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'utilisation de l'item.")
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez utilisé %d %s sur <@%s>.", quantity, itemName, targetID))
	if err != nil {
		log.Println("Erreur lors de l'envoi du message:", err)
	}
}
