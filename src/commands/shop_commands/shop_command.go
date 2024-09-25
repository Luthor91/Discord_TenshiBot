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

// Constantes pour les noms des packs d'XP
const (
	PetitPackXP = "smolpack"
	MoyenPackXP = "pack"
	GrandPackXP = "bigpack"
	Timeout     = "timeout"
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

	UserDiscordID := m.Author.ID

	// Récupère les informations de l'utilisateur
	userService := services.NewUserService()
	userMoney, _ := userService.GetMoney(UserDiscordID)
	userXP, _ := userService.GetExperience(UserDiscordID)

	// Récupérer les items de la base de données
	shopService := services.NewShopService()
	options, err := shopService.GetShopItems()
	if err != nil {
		log.Println("Erreur lors de la récupération des items:", err)
		return
	}

	// Prépare le message du shop
	messageContent := "**Bienvenue dans le shop !**\n\n"
	for _, option := range options {
		messageContent += fmt.Sprintf("%s **Acheter %s pour %.2f money**\n", option.Emoji, option.Name, option.Price)
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

		// Vérifier si le cooldown existe
		exists, err := shopService.IsUserShopCooldownExists(UserDiscordID, selectedOption.ID)
		if err != nil {
			log.Println("Erreur lors de la vérification du cooldown:", err)
			return
		}

		// Si le cooldown existe, refuse l'achat
		if exists {
			_, err := s.ChannelMessageSend(m.ChannelID, "Vous ne pouvez pas acheter cet article en raison d'un cooldown actif.")
			if err != nil {
				log.Println("Erreur lors de l'envoi du message:", err)
			}
			return
		}

		// Vérifier l'argent de l'utilisateur et appliquer l'action
		if userMoney >= int(selectedOption.Price) {
			var err error
			user, err := userService.GetUserByDiscordID(UserDiscordID)
			if err != nil {
				log.Println("Erreur lors de la récupération de l'utilisateur:", err)
				return
			}

			// Appliquer les effets de l'achat
			switch selectedOption.Name {
			case PetitPackXP:
				userService.AddExperience(user.UserDiscordID, 50)
				userService.AddMoney(user.UserDiscordID, -100)
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté un %s pour 100 money.", PetitPackXP))
			case MoyenPackXP:
				userService.AddExperience(user.UserDiscordID, 500)
				userService.AddMoney(user.UserDiscordID, -1050)
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté un %s pour 1000 money.", MoyenPackXP))
			case GrandPackXP:
				userService.AddExperience(user.UserDiscordID, 5000)
				userService.AddMoney(user.UserDiscordID, -11000)
				_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté un %s pour 5000 money.", GrandPackXP))
			case Timeout:
				services.NewItemService().AddItem(UserDiscordID, "timeout", 1)
				userService.AddMoney(user.UserDiscordID, -5000)
				_, _ = s.ChannelMessageSend(m.ChannelID, "Vous avez acheté un timeout de 5 minutes pour 5000 money.")
			}

			// Initialiser le cooldown après l'achat
			now := time.Now()
			cooldown := &models.UserShopCooldown{
				UserDiscordID: UserDiscordID,
				ItemID:        selectedOption.ID,
				NextPurchase:  now.Add(time.Duration(selectedOption.Cooldown) * time.Second),
			}
			err = shopService.SetUserShopCooldown(cooldown.UserDiscordID, cooldown.ItemID, cooldown.NextPurchase)
			if err != nil {
				log.Println("Erreur lors de la mise à jour du cooldown:", err)
				return
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
			if err != nil {
				log.Println("Erreur lors de l'envoi du message:", err)
			}
		}
	})
}
