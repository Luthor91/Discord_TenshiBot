package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TimeoutCommand est la commande qui applique un timeout
func TimeoutCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Utilisation: !timeout @mention durée[s/m/h/d]")
		return
	}

	// Vérifier s'il y a une mention d'utilisateur
	if len(m.Mentions) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur.")
		return
	}

	// Récupérer l'ID de l'utilisateur mentionné
	userID := m.Mentions[0].ID
	durationStr := args[2]

	// Parser la durée
	duration, err := parseDuration(durationStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors du parsing de la durée : %v", err))
		return
	}

	// Appliquer le timeout
	err = timeoutUser(s, m.GuildID, userID, duration)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la mise en timeout de l'utilisateur : %v", err))
		return
	}

	// Informer du succès de la commande
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'utilisateur <@%s> a été mis en timeout pour %v.", userID, duration))
}

// timeoutUser met un utilisateur en timeout pour une durée donnée
func timeoutUser(s *discordgo.Session, guildID, userID string, duration time.Duration) error {
	timeoutUntil := time.Now().Add(duration)
	err := s.GuildMemberTimeout(guildID, userID, &timeoutUntil)
	if err != nil {
		return err
	}
	return nil
}

// parseDuration parse la durée au format '10s', '5m', '2h', '1d'
func parseDuration(durationStr string) (time.Duration, error) {
	if len(durationStr) < 2 {
		return 0, fmt.Errorf("durée invalide")
	}

	unit := durationStr[len(durationStr)-1]
	valueStr := durationStr[:len(durationStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("valeur numérique invalide")
	}

	switch unit {
	case 's': // secondes
		return time.Duration(value) * time.Second, nil
	case 'm': // minutes
		return time.Duration(value) * time.Minute, nil
	case 'h': // heures
		return time.Duration(value) * time.Hour, nil
	case 'd': // jours
		return time.Duration(value) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("unité de temps invalide")
	}
}
