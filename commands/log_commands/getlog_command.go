package log_commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Commande pour récupérer les X derniers logs
func GetLogsCommand(s *discordgo.Session, m *discordgo.MessageCreate, service *services.LogService) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if !isMod {
		return
	}

	// Vérifier si le message commence par la commande de récupération des logs utilisateur
	command := fmt.Sprintf("%slogs", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments après la commande (nombre de logs à récupérer et option verbose)
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer combien de logs vous voulez récupérer.")
		return
	}

	// Extraire la limite
	limit, err := strconv.Atoi(args[1])
	if err != nil || limit <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
		return
	}

	// Récupérer les derniers logs
	logs, err := service.GetLastLogs(limit)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs.")
		log.Printf("Erreur lors de la récupération des logs : %v", err)
		return
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
