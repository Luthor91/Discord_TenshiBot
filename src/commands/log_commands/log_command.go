package log_commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/models"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Commande pour récupérer les logs
func LogsCommand(s *discordgo.Session, m *discordgo.MessageCreate, service *services.LogService) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if !isMod {
		return
	}

	// Vérifier si le message commence par la commande de récupération des logs
	command := fmt.Sprintf("%slogs", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments après la commande
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer combien de logs vous voulez récupérer.")
		return
	}

	// Définir des variables pour l'utilisateur et la limite
	var userID string
	limit := 0
	var err error

	// Vérifier si l'argument "-n" est présent pour récupérer les logs d'un utilisateur spécifique
	if args[1] == "-n" && len(args) >= 4 {
		userID = args[2]
		// Nettoyer l'ID de l'utilisateur mentionné
		if strings.HasPrefix(userID, "<@") {
			userID = strings.Trim(userID, "<@!>")
		}
		// Extraire la limite (nombre de logs à récupérer)
		limit, err = strconv.Atoi(args[3])
		if err != nil || limit <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
			return
		}
	} else {
		// Pas de "-n", donc récupérer les derniers logs globaux
		limit, err = strconv.Atoi(args[1])
		if err != nil || limit <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
			return
		}
	}

	// Récupérer les logs selon le contexte (utilisateur ou logs globaux)
	var logs []models.Log
	if userID != "" {
		// Récupérer les logs de l'utilisateur spécifié
		logs, err = service.GetLogsByUser(userID, limit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs de l'utilisateur.")
			log.Printf("Erreur lors de la récupération des logs de l'utilisateur %s : %v", userID, err)
			return
		}
	} else {
		// Récupérer les derniers logs globaux
		logs, err = service.GetLastLogs(limit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs.")
			log.Printf("Erreur lors de la récupération des logs : %v", err)
			return
		}
	}

	// Vérifier s'il y a des logs disponibles
	if len(logs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Aucun log trouvé.")
		return
	}

	// Envoyer les logs dans le canal
	for _, logEntry := range logs {
		message := fmt.Sprintf("`[%s]` **%s** dans **%s**: %s", logEntry.Timestamp.Format(time.RFC822), logEntry.Username, logEntry.ChannelName, logEntry.Message)
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
