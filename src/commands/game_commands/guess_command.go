package game_commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// GuessCommand devine un nombre
func GuessCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sguess", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?guess <number>")
		return
	}

	guess, err := strconv.Atoi(parts[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un nombre valide.")
		return
	}

	// Générer un nombre aléatoire entre 1 et 10
	target := rand.Intn(10) + 1 // 1 à 10

	// Vérifier si le joueur a gagné
	if guess == target {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bravo ! Vous avez deviné le bon nombre : %d.", target))
		// Récompense de 10 money
		userService := services.NewUserService()
		user, _ := userService.GetUserByDiscordID(m.Author.ID)
		userService.AddMoney(user.UserDiscordID, 10)
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Dommage ! Le nombre était %d.", target))
	}
}
