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
func LogsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	// Définir des variables pour l'utilisateur, le salon et la limite
	var userID string
	var channelID string
	limit := 0
	var err error

	// Vérifier si l'argument "-n" est présent pour récupérer les logs d'un utilisateur spécifique
	for i, arg := range args {
		switch arg {
		case "-n":
			if i+1 < len(args) {
				userID = args[i+1]
				// Nettoyer l'ID de l'utilisateur mentionné
				if strings.HasPrefix(userID, "<@") {
					userID = strings.Trim(userID, "<@!>")
				}
			}
		case "-c":
			if i+1 < len(args) {
				channelID = args[i+1]
				// Vérifier si c'est une mention de salon
				if strings.HasPrefix(channelID, "<#") {
					channelID = strings.Trim(channelID, "<#>")
				}
			}
		}
	}

	// Extraire la limite (nombre de logs à récupérer) depuis les arguments
	if userID != "" {
		// Si un utilisateur est spécifié, la limite doit être le prochain argument
		if len(args) >= 4 {
			limit, err = strconv.Atoi(args[len(args)-1]) // Dernier argument après les options
			if err != nil || limit <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
				return
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer combien de logs vous voulez récupérer.")
			return
		}
	} else {
		// Pas d'utilisateur spécifié, la limite est le premier argument
		limit, err = strconv.Atoi(args[1])
		if err != nil || limit <= 0 {
			s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
			return
		}
	}

	// Récupérer les logs selon le contexte (utilisateur, salon, ou logs globaux)
	var logs []models.Log
	logService := services.NewLogService()
	if userID != "" && channelID != "" {
		// Récupérer les logs de l'utilisateur spécifié dans le salon spécifié
		logs, err = logService.GetLogsByUserAndChannel(userID, channelID, limit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs de l'utilisateur dans le salon spécifié.")
			log.Printf("Erreur lors de la récupération des logs de l'utilisateur %s dans le salon %s : %v", userID, channelID, err)
			return
		}
	} else if userID != "" {
		// Récupérer les logs de l'utilisateur spécifié
		logs, err = logService.GetLogsByUser(userID, limit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs de l'utilisateur.")
			log.Printf("Erreur lors de la récupération des logs de l'utilisateur %s : %v", userID, err)
			return
		}
	} else if channelID != "" {
		// Récupérer les logs du salon spécifié
		logs, err = logService.GetLogsByChannel(channelID, limit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs du salon.")
			log.Printf("Erreur lors de la récupération des logs du salon %s : %v", channelID, err)
			return
		}
	} else {
		// Récupérer les derniers logs globaux
		logs, err = logService.GetLastLogs(limit)
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
