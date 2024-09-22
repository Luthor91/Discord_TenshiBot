package shop_commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// BuyCommand permet d'acheter un item directement en spécifiant son nom et sa quantité
func BuyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sbuy", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	args := strings.Fields(m.Content[len(command):])
	if len(args) < 2 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Usage : ?buy <item_name> <quantity>")
		return
	}

	itemName := args[0]
	quantity, err := strconv.Atoi(args[1])
	if err != nil || quantity <= 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier une quantité valide.")
		return
	}

	userID := m.Author.ID
	userMoney, err := services.NewUserService(controllers.NewUserController()).GetMoney(userID)
	if err != nil {
		log.Println("Erreur lors de la récupération de l'argent de l'utilisateur:", err)
		return
	}

	// Récupérer l'item par son nom
	item, err := services.NewShopService().GetShopItemByName(itemName)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "L'item spécifié n'existe pas.")
		return
	}

	totalCost := int(item.Price) * quantity
	if userMoney < totalCost {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
		return
	}

	// Appliquer l'achat
	err = services.NewUserService(controllers.NewUserController()).UpdateMoney(userID, -totalCost)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Erreur lors de la mise à jour de votre argent.")
		log.Println("Erreur lors de la mise à jour de l'argent:", err)
		return
	}

	// Ajouter l'item à l'utilisateur
	err = services.NewItemService().AddItem(userID, item.Name, quantity)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout de l'item.")
		log.Println("Erreur lors de l'ajout de l'item:", err)
		return
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté %d %s pour %d money.", quantity, item.Name, totalCost))
}
