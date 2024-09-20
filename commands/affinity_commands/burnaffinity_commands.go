package affinity_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// BurnAffinity retire un montant d'affinité à l'utilisateur
func BurnAffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commande de base : ?burnaffinity <quantité>
	prefix := fmt.Sprintf("%sburnaffinity", config.AppConfig.BotPrefix)
	if strings.HasPrefix(m.Content, prefix) {
		args := strings.Fields(m.Content)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?burnaffinity <quantité>")
			return
		}

		// Convertir la quantité en entier
		affinityAmount, err := strconv.Atoi(args[1])
		if err != nil || affinityAmount < 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'affinité.")
			return
		}

		// Récupérer l'affinité de l'utilisateur
		user, exists := services.NewAffinityService(s).GetUserAffinity(m.Author.ID)
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé.")
			return
		}

		// Vérifier si l'utilisateur a assez d'affinité à retirer
		if user.Affinity < affinityAmount {
			s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'affinité pour retirer ce montant.")
			return
		}

		// Définir l'affinité pour l'utilisateur
		newAffinity := user.Affinity - affinityAmount
		err = services.NewUserService(controllers.NewUserController()).SetAffinity(m.Author.ID, newAffinity)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'affinité.")
			return
		}
	}
}
