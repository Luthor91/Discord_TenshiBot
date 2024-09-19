package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

const (
	cooldownDuration = time.Hour       // 1 heure de cooldown pour les options
	timeoutDuration  = 5 * time.Minute // 5 minutes de timeout
)

// ShopItem représente un article dans le magasin avec son nom, son prix, son cooldown, et son emoji
type ShopItem struct {
	ID       uint
	Name     string
	Price    float64
	Cooldown int
	Emoji    string // Emoji pour la réaction Discord
	Action   func(userID string) string
}

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
	userMoney, _ := services.GetUserMoney(userID)
	userXP, _ := services.GetExperience(userID)

	// Options du shop
	options := []ShopItem{
		{
			ID:       1,
			Name:     "XP Pack 50",
			Price:    100,
			Cooldown: int(cooldownDuration.Seconds()),
			Emoji:    "1️⃣",
			Action: func(userID string) string {
				services.AddExperience(userID, 50)
				services.AddMoney(userID, -100)
				return "Vous avez acheté 50 XP pour 100 money."
			},
		},
		{
			ID:       2,
			Name:     "XP Pack 500",
			Price:    1000,
			Cooldown: int(cooldownDuration.Seconds()),
			Emoji:    "2️⃣",
			Action: func(userID string) string {
				services.AddExperience(userID, 500)
				services.AddMoney(userID, -1000)
				return "Vous avez acheté 500 XP pour 1000 money."
			},
		},
		{
			ID:       3,
			Name:     "XP pour Money",
			Price:    float64(userMoney) * 0.15,
			Cooldown: int(cooldownDuration.Seconds()),
			Emoji:    "3️⃣",
			Action: func(userID string) string {
				xpToAdd := int(float64(userXP) * 0.20)
				cost := float64(userMoney) * 0.15
				services.AddExperience(userID, xpToAdd)
				services.AddMoney(userID, -int(cost))
				return fmt.Sprintf("Vous avez acheté %d XP pour %d money.", xpToAdd, int(cost))
			},
		},
		{
			ID:       4,
			Name:     "Timeout",
			Price:    5000,
			Cooldown: int(timeoutDuration.Seconds()),
			Emoji:    "4️⃣",
			Action: func(userID string) string {
				services.AddItem(userID, "timeout", 1)
				services.AddMoney(userID, -5000)
				return "Vous avez acheté un timeout de 5 minutes pour 5000 money."
			},
		},
	}

	// Prépare le message du shop
	messageContent := fmt.Sprintf(
		"**Bienvenue dans le shop !**\n\n"+
			"1️⃣ **Acheter 50 XP pour 100 money**\n"+
			"2️⃣ **Acheter 500 XP pour 1000 money**\n"+
			"3️⃣ **Acheter %d XP pour %d money**\n"+
			"4️⃣ **Acheter un timeout de 5 minutes pour 5000 money**\n\n"+
			"Votre solde actuel : %d money\n"+
			"Votre XP actuel : %d",
		int(float64(userXP)*0.20), int(float64(userMoney)*0.15),
		userMoney, userXP,
	)

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

		// Recharger les cooldowns au cas où ils ont changé
		cooldowns, err := services.LoadUserShopCooldowns(userID)
		if err != nil {
			log.Println("Erreur lors du rechargement des cooldowns:", err)
			return
		}

		// Trouver l'option correspondant à l'emoji
		for _, option := range options {
			if r.Emoji.Name == option.Emoji {
				// Vérifier le cooldown
				userCooldown, ok := cooldowns[option.ID]
				if !ok {
					userCooldown = models.UserShopCooldown{NextPurchase: time.Time{}}
				}
				now := time.Now()

				if now.Sub(userCooldown.NextPurchase) < time.Duration(option.Cooldown)*time.Second {
					remaining := time.Duration(option.Cooldown)*time.Second - now.Sub(userCooldown.NextPurchase)
					_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'option est en cooldown. Temps restant: %s", remaining))
					if err != nil {
						log.Println("Erreur lors de l'envoi du message:", err)
					}
					return
				}

				// Vérifier l'argent et appliquer l'action
				if userMoney >= int(option.Price) {
					response := option.Action(userID)
					_, err := s.ChannelMessageSend(m.ChannelID, response)
					if err != nil {
						log.Println("Erreur lors de l'envoi du message:", err)
					}

					// Mettre à jour le cooldown
					nextPurchase := time.Now().Add(time.Duration(option.Cooldown) * time.Second)
					_, err = controllers.UpdateUserShopCooldown(userID, option.ID, nextPurchase)
					if err != nil {
						log.Println("Erreur lors de la mise à jour du cooldown:", err)
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
					if err != nil {
						log.Println("Erreur lors de l'envoi du message:", err)
					}
				}
				return
			}
		}
	})
}
