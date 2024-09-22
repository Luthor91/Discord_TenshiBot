package game_commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// ShifumiCommand joue à Pierre-Papier-Ciseaux
func ShifumiCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sshifumi", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	parts := strings.Fields(m.Content)
	if len(parts) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: ?shifumi <pierre | feuille | ciseaux>")
		return
	}

	playerChoice := strings.ToLower(parts[1])
	choices := map[string]string{"pierre": "rock", "feuille": "paper", "ciseaux": "scissors"}

	// Vérifier que le choix de l'utilisateur est valide
	if _, exists := choices[playerChoice]; !exists {
		s.ChannelMessageSend(m.ChannelID, "Veuillez choisir parmi : pierre, feuille, ou ciseaux.")
		return
	}

	// Choix du bot
	botChoice := rand.Intn(3) // 0: pierre, 1: feuille, 2: ciseaux
	var botChoiceStr string
	switch botChoice {
	case 0:
		botChoiceStr = "pierre"
	case 1:
		botChoiceStr = "feuille"
	case 2:
		botChoiceStr = "ciseaux"
	}

	// Déterminer le résultat
	var result string
	if playerChoice == botChoiceStr {
		result = "C'est une égalité !"
	} else if (playerChoice == "pierre" && botChoiceStr == "ciseaux") ||
		(playerChoice == "feuille" && botChoiceStr == "pierre") ||
		(playerChoice == "ciseaux" && botChoiceStr == "feuille") {
		result = "Vous avez gagné !"
		// Récompense de 10 money
		userService := services.NewUserService(controllers.NewUserController())
		user, _ := userService.GetUserByDiscordID(m.Author.ID)
		userService.AddMoney(user, 10)
	} else {
		result = "Vous avez perdu !"
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Bot a choisi : %s. %s", botChoiceStr, result))
}
