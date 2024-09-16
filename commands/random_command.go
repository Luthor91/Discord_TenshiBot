package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

// RandomCommand génère un nombre aléatoire entre deux bornes spécifiées
func RandomCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%srandom", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if strings.HasPrefix(m.Content, command) {
		// Extraire les arguments après la commande (bornes pour le nombre aléatoire)
		args := strings.Fields(m.Content)
		if len(args) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Merci de spécifier les bornes pour le tirage aléatoire.")
			return
		}

		// Convertir les bornes en entiers
		min, err := strconv.Atoi(args[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une borne inférieure valide.")
			return
		}
		max, err := strconv.Atoi(args[2])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une borne supérieure valide.")
			return
		}

		// Initialiser le générateur de nombres aléatoires
		rand.Seed(time.Now().UnixNano())
		randomNumber := rand.Intn(max-min+1) + min

		// Envoyer le nombre aléatoire généré
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Nombre aléatoire entre %d et %d : %d", min, max, randomNumber))
	}
}
