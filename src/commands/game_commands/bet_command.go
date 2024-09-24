package game_commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// BetCommand permet à un utilisateur de parier une somme d'argent
func BetCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorer les messages du bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sbet", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer les arguments de la commande
	parts := strings.Fields(m.Content)
	if len(parts) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?bet <montant>")
		return
	}

	// Récupérer le montant à parier
	amountStr := parts[1]
	amount, err := strconv.Atoi(amountStr)
	if err != nil || amount <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide.")
		return
	}

	// Créer une instance de UserService
	userService := services.NewUserService()

	// Récupérer les informations de l'utilisateur
	userID := m.Author.ID
	userMoney, err := userService.GetMoney(userID) // Utiliser UserService pour récupérer le solde
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de votre solde.")
		return
	}

	// Vérifier si l'utilisateur a suffisamment d'argent
	if userMoney < amount {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour parier ce montant.")
		return
	}

	// Initialiser le générateur de nombres aléatoires
	rand.Seed(time.Now().UnixNano())
	win := rand.Intn(2) // 0 ou 1, 50% de chances

	if win == 1 {
		// L'utilisateur gagne
		newBalance := userMoney + amount
		err := userService.SetMoney(userID, newBalance) // Utiliser UserService pour mettre à jour le solde
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la mise à jour de votre solde.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🎉 Vous avez gagné ! Votre nouveau solde est de %d.", newBalance))
	} else {
		// L'utilisateur perd
		newBalance := userMoney - amount
		err := userService.SetMoney(userID, newBalance) // Utiliser UserService pour mettre à jour le solde
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la mise à jour de votre solde.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("😢 Vous avez perdu. Votre nouveau solde est de %d.", newBalance))
	}
}
