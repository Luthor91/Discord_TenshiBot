package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/bwmarrin/discordgo"
)

const (
	// Identifiants des réactions pour sélectionner les options
	option1Emoji     = "1️⃣"
	option2Emoji     = "2️⃣"
	option3Emoji     = "3️⃣"
	cooldownDuration = time.Hour // 1 heure de cooldown
)

// ShopCommand affiche le magasin avec des options pour dépenser de la money pour de l'xp
func ShopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sshop", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	userID := m.Author.ID

	// Charger les cooldowns depuis le fichier
	cooldowns, err := features.LoadCooldowns()
	if err != nil {
		log.Println("Erreur lors du chargement des cooldowns:", err)
		return
	}

	// Vérifier si l'utilisateur a déjà interagi avec le shop
	userCooldown, exists := cooldowns[userID]
	if !exists {
		// L'utilisateur n'a jamais interagi avec le shop, initialiser ses cooldowns
		userCooldown = features.UserCooldowns{
			Option1: time.Time{}, // Pas encore utilisé
			Option2: time.Time{},
			Option3: time.Time{},
		}
		cooldowns[userID] = userCooldown

		// Sauvegarder les cooldowns mis à jour dans le fichier
		err = features.SaveCooldowns(cooldowns)
		if err != nil {
			log.Println("Erreur lors de la sauvegarde des cooldowns:", err)
			return
		}
	}

	// Récupère l'argent et l'expérience de l'utilisateur
	userMoney := features.GetUserMoney(userID)
	userXP, _ := features.GetExperience(userID)

	// Options du magasin
	option1XP := 50             // Nombre d'XP pour l'option 1
	option1Cost := 100          // Coût en money pour l'option 1
	option2XP := 500            // Nombre d'XP pour l'option 2
	option2Cost := 1000         // Coût en money pour l'option 2
	option3PercentageXP := 20   // Pourcentage de l'XP de l'utilisateur pour l'option 3
	option3CostPercentage := 15 // Pourcentage du montant de money pour l'option 3

	// Calcul des valeurs pour les options du shop
	option3XP := int(float64(userXP) * float64(option3PercentageXP) / 100)
	option3Cost := int(float64(userMoney) * float64(option3CostPercentage) / 100)

	// Prépare le message
	messageContent := fmt.Sprintf(
		"**Bienvenue dans le shop !**\n\n"+
			"1️⃣ **Acheter %d XP pour %d money**\n"+
			"2️⃣ **Acheter %d XP pour %d money**\n"+
			"3️⃣ **Acheter %d XP pour %d money**\n\n"+
			"Votre solde actuel : %d money\n"+
			"Votre XP actuel : %d",
		option1XP, option1Cost,
		option2XP, option2Cost,
		option3XP, option3Cost,
		userMoney, userXP,
	)

	// Envoi du message et ajout des réactions
	msg, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil {
		log.Println("Erreur lors de l'envoi du message:", err)
		return
	}

	// Ajouter des réactions au message
	for _, react := range []string{option1Emoji, option2Emoji, option3Emoji} {
		err = s.MessageReactionAdd(m.ChannelID, msg.ID, react)
		if err != nil {
			log.Println("Erreur lors de l'ajout de la réaction:", err)
		}
	}

	// Gestion des réactions
	s.AddHandler(func(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
		if r.MessageID != msg.ID || r.UserID == s.State.User.ID {
			return
		}

		// Charger les cooldowns depuis le fichier avant de vérifier le cooldown pour la réaction
		cooldowns, err := features.LoadCooldowns()
		if err != nil {
			log.Println("Erreur lors du chargement des cooldowns:", err)
			return
		}

		userCooldown := cooldowns[userID]
		now := time.Now()

		switch r.Emoji.Name {
		case option1Emoji:
			// Vérifie le cooldown pour l'option 1
			if now.Sub(userCooldown.Option1) < cooldownDuration {
				remaining := cooldownDuration - now.Sub(userCooldown.Option1)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'option 1 est en cooldown. Temps restant: %s", remaining))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				return
			}

			if userMoney >= option1Cost {
				features.AddExperience(userID, option1XP)
				features.AddMoney(userID, -option1Cost)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté %d XP pour %d money.", option1XP, option1Cost))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				// Met à jour le cooldown pour l'option 1
				userCooldown.Option1 = time.Now()
				cooldowns[userID] = userCooldown
				err = features.SaveCooldowns(cooldowns)
				if err != nil {
					log.Println("Erreur lors de la sauvegarde des cooldowns:", err)
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
			}

		case option2Emoji:
			// Vérifie le cooldown pour l'option 2
			if now.Sub(userCooldown.Option2) < cooldownDuration {
				remaining := cooldownDuration - now.Sub(userCooldown.Option2)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'option 2 est en cooldown. Temps restant: %s", remaining))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				return
			}

			if userMoney >= option2Cost {
				features.AddExperience(userID, option2XP)
				features.AddMoney(userID, -option2Cost)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté %d XP pour %d money.", option2XP, option2Cost))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				// Met à jour le cooldown pour l'option 2
				userCooldown.Option2 = time.Now()
				cooldowns[userID] = userCooldown
				err = features.SaveCooldowns(cooldowns)
				if err != nil {
					log.Println("Erreur lors de la sauvegarde des cooldowns:", err)
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
			}

		case option3Emoji:
			// Vérifie le cooldown pour l'option 3
			if now.Sub(userCooldown.Option3) < cooldownDuration {
				remaining := cooldownDuration - now.Sub(userCooldown.Option3)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'option 3 est en cooldown. Temps restant: %s", remaining))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				return
			}

			if userMoney >= option3Cost {
				features.AddExperience(userID, option3XP)
				features.AddMoney(userID, -option3Cost)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez acheté %d XP pour %d money.", option3XP, option3Cost))
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
				// Met à jour le cooldown pour l'option 3
				userCooldown.Option3 = time.Now()
				cooldowns[userID] = userCooldown
				err = features.SaveCooldowns(cooldowns)
				if err != nil {
					log.Println("Erreur lors de la sauvegarde des cooldowns:", err)
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez de money.")
				if err != nil {
					log.Println("Erreur lors de l'envoi du message:", err)
				}
			}
		}
	})
}
