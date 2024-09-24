package shop_commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/models"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// ShopCommand affiche le magasin avec des options pour dépenser de la money pour de l'XP ou d'autres items
func ShopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sshop", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	userID := m.Author.ID

	// Récupère les informations de l'utilisateur
	userService := services.NewUserService()
	userMoney, _ := userService.GetMoney(userID)
	userXP, _ := userService.GetExperience(userID)

	// Récupérer les items de la base de données
	shopService := services.NewShopService()
	options, err := shopService.GetShopItems()
	if err != nil {
		log.Println("Erreur lors de la récupération des items:", err)
		return
	}

	// Prépare le message du shop
	messageContent := "**Bienvenue dans le shop !**\n\n"
	for i, option := range options {
		if i == 2 {
			xpToAdd := int(float64(userXP) * 0.10)
			cost := float64(userMoney) * 0.30
			messageContent += fmt.Sprintf("%s **Acheter %d %s pour %.2f money**\n", option.Emoji, xpToAdd, option.Name, cost)
		} else {
			messageContent += fmt.Sprintf("%s **Acheter %s pour %.2f money**\n", option.Emoji, option.Name, option.Price)
		}
	}

	messageContent += fmt.Sprintf("\nVotre solde actuel : %d money\nVotre XP actuel : %d", userMoney, userXP)

	// Envoi du message et ajout des réactions
	msg, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil {
		log.Println("Erreur lors de l'envoi du message:", err)
		return
	}

	// Ajouter des réactions au message
	for _, option := range options {
		err = s.MessageReactionAdd(m.ChannelID, msg.ID, option.Emoji)
		if err != nil {
			log.Println("Erreur lors de l'ajout de la réaction:", err)
		}
	}

	// Gestion des réactions
	s.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.MessageID != msg.ID || r.UserID == s.State.User.ID {
			return
		}

		// Trouver l'option correspondant à l'emoji
		var selectedOption models.ShopItem
		for _, option := range options {
			if r.Emoji.Name == option.Emoji {
				selectedOption = option
				break
			}
		}

		// Recharger les cooldowns
		cooldownService := services.NewShopService()
		cooldown, err := cooldownService.GetUserShopCooldown(userID, selectedOption.ID)
		if err != nil {
			log.Println("Erreur lors du rechargement des cooldowns:", err)
			return
		}

		now := time.Now()

		// Si le cooldown est nil, l'initialiser
		if cooldown == nil {
			cooldown = &models.UserShopCooldown{UserDiscordID: userID, ItemID: selectedOption.ID, NextPurchase: time.Time{}}
		}

		// Vérifier le temps de cooldown
		if now.Sub(cooldown.NextPurchase) < time.Duration(selectedOption.Cooldown)*time.Second {
			remaining := time.Duration(selectedOption.Cooldown)*time.Second - now.Sub(cooldown.NextPurchase)
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'option est en cooldown. Temps restant: %s", remaining))
			if err != nil {
				log.Println("Erreur lors de l'envoi du message:", err)
			}
			return
		}

		// Vérifier l'argent et appliquer l'action
		if userMoney >= int(selectedOption.Price) {
			var err error
			user, err := userService.GetUserByDiscordID(userID)
			if err != nil {
				log.Println("Erreur lors de la récupération de l'utilisateur:", err)
				return
			}

			// Appliquer les effets de l'achat
			switch selectedOption.Name {
			case "50 XP":
				userService.AddExperience(user, 50)
				userService.AddMoney(user, -100)
				_, _ = s.ChannelMessageSend(m.ChannelID, "Vous avez acheté 50 XP pour 100 money.")
			case "500 XP":
				userService.AddExperience(user, 500)
				userService.AddMoney(user, -1000)
				_, _ = s.ChannelMessageSend(m.ChannelID, "Vous avez acheté 500 XP pour 1000 money.")
			case "XP":
				xpToAdd := int(float64(userXP) * 0.10)
				cost := float64(userMoney) * 0.30
				userService.AddExperience(user, xpToAdd)
				userService.AddMoney(user, -int(cost))
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté %d XP pour %.2f money.", xpToAdd, cost))
			case "Timeout":
				services.NewItemService().AddItem(userID, "timeout", 1)
				userService.AddMoney(user, -5000)
				_, _ = s.ChannelMessageSend(m.ChannelID, "Vous avez acheté un timeout de 5 minutes pour 5000 money.")
			}

			// Mettre à jour le cooldown après l'achat
			cooldown.NextPurchase = now.Add(time.Duration(selectedOption.Cooldown) * time.Second)
			err = cooldownService.SetUserShopCooldown(userID, selectedOption.ID, cooldown.NextPurchase)
			if err != nil {
				log.Println("Erreur lors de la mise à jour du cooldown:", err)
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
			if err != nil {
				log.Println("Erreur lors de l'envoi du message:", err)
			}
		}
	})
}
