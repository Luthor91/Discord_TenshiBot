package item_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// ItemCommand gère les différentes actions sur les items
func ItemCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sitem", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments de la commande
	args := strings.Fields(m.Content[len(command):])
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?item [-u|-r|-m|-g] [-n <mention>] <item_name> <quantity>")
		return
	}

	action := args[0]
	var itemName string
	var quantity int
	var targetID string
	var err error

	// Analyser l'option -n pour la mention
	for i, arg := range args {
		if arg == "-n" && i+1 < len(args) {
			targetMention := args[i+1]
			targetID = strings.TrimPrefix(targetMention, "<@!")
			targetID = strings.TrimSuffix(targetID, ">")
		} else if arg == "-u" || arg == "-r" || arg == "-g" {
			action = arg
		} else if i == len(args)-2 {
			itemName = arg
		} else if i == len(args)-1 {
			quantity, err = strconv.Atoi(arg)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Quantité invalide.")
				return
			}
		}
	}

	// Action : Utiliser un item
	if action == "-u" {
		err = services.NewItemService().UseItem(m.Author.ID, targetID, itemName, quantity)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'utilisation de l'item.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez utilisé %d %s sur <@%s>.", quantity, itemName, targetID))

		// Action : Jeter un item
	} else if action == "-r" {
		err = services.NewItemService().RemoveItem(m.Author.ID, itemName, quantity)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait de l'item.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez jeté %d %s.", quantity, itemName))

		// Action : Lister l'inventaire
	} else if action == "-m" {
		items, err := services.NewItemService().GetUserItems(m.Author.ID)
		if err != nil || len(items) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Votre inventaire est vide ou une erreur est survenue.")
			return
		}
		messageContent := "**Votre inventaire :**\n"
		for _, item := range items {
			messageContent += fmt.Sprintf("%s: %d\n", item.Name, item.Quantity)
		}
		s.ChannelMessageSend(m.ChannelID, messageContent)

		// Action : Donner un item
	} else if action == "-g" {
		err = services.NewItemService().GiveItem(m.Author.ID, targetID, itemName, quantity)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du transfert de l'item.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez donné %d %s à <@%s>.", quantity, itemName, targetID))
	}
}
