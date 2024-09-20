package item_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// BurnItem retire un item de l'inventaire de l'utilisateur
func BurnItem(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?burnitem <nom_item> <quantité>
	prefix := fmt.Sprintf("%sburnitem", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?burnitem <nom_item> <quantité>")
			return
		}

		itemName := args[1]
		itemAmount, err := strconv.Atoi(args[2])
		if err != nil || itemAmount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'item.")
			return
		}

		// Récupérer les items de l'utilisateur
		items, err := services.NewItemService().GetUserItems(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des items.")
			return
		}

		// Vérifier si l'utilisateur a assez d'items à retirer
		itemFound := false
		for _, item := range items {
			if item.Name == itemName {
				itemFound = true
				if item.Quantity < itemAmount {
					s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'items pour retirer ce montant.")
					return
				}
				break
			}
		}

		if !itemFound {
			s.ChannelMessageSend(m.ChannelID, "Item non trouvé dans votre inventaire.")
			return
		}

		// Retirer l'item
		err = services.NewItemService().RemoveItem(m.Author.ID, itemName, itemAmount)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'item.")
			return
		}

		// Confirmation à l'utilisateur
		response := fmt.Sprintf("Vous avez retiré %d de l'item '%s'.", itemAmount, itemName)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
