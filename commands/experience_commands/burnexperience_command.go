package experience_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// BurnExperienceCommand retire un montant d'XP à l'utilisateur
func BurnExperienceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?burnxp <quantité>
	prefix := fmt.Sprintf("%sburnxp", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?burnxp <quantité>")
			return
		}

		// Convertir la quantité en entier
		xpAmount, err := strconv.Atoi(args[1])
		if err != nil || xpAmount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'XP.")
			return
		}

		// Récupérer l'XP de l'utilisateur
		amount, exists := services.NewUserService(controllers.NewUserController()).GetExperience(m.Author.ID)
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé.")
			return
		}

		// Vérifier si l'utilisateur a assez d'XP à retirer
		if amount < xpAmount {
			s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'XP pour retirer ce montant.")
			return
		}

		// Définir l'XP pour l'utilisateur
		newXP := amount - xpAmount
		err = services.NewUserService(controllers.NewUserController()).SetExperience(m.Author.ID, newXP)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'XP.")
			return
		}

		// Confirmation à l'utilisateur
		response := fmt.Sprintf("Vous avez retiré %d d'XP. Votre nouvel XP est de %d.", xpAmount, newXP)
		s.ChannelMessageSend(m.ChannelID, response)
	}
}
