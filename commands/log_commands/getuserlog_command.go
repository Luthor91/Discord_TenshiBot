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

// Commande pour récupérer les logs d'un utilisateur spécifique
func GetUserLogsCommand(s *discordgo.Session, m *discordgo.MessageCreate, service *services.LogService) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	// Vérifier si le message commence par la commande de récupération des logs utilisateur
	command := fmt.Sprintf("%suserlogs", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments après la commande (utilisateur et nombre de logs)
	args := strings.Fields(m.Content)
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur et indiquer combien de logs vous voulez récupérer.")
		return
	}

	// Extraire l'ID de l'utilisateur mentionné
	userID := args[1]
	if !strings.HasPrefix(userID, "<@") {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide.")
		return
	}
	// Nettoyer l'ID de l'utilisateur mentionné
	userID = strings.Trim(userID, "<@!>")

	// Extraire la limite (nombre de logs à récupérer)
	limit, err := strconv.Atoi(args[2])
	if err != nil || limit <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez indiquer un nombre valide de logs.")
		return
	}

	// Récupérer les logs de l'utilisateur spécifié
	logs, err := service.GetLogsByUser(userID, limit)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des logs de l'utilisateur.")
		log.Printf("Erreur lors de la récupération des logs de l'utilisateur %s : %v", userID, err)
		return
	}

	// Vérifier s'il y a des logs disponibles pour l'utilisateur
	if len(logs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Aucun log trouvé pour cet utilisateur.")
		return
	}

	// Envoyer les logs dans le canal
	for _, logEntry := range logs {
		message := fmt.Sprintf("`[%s]` **%s** dans **%s**: %s", logEntry.Timestamp.Format(time.RFC822), logEntry.Username, logEntry.ChannelName, logEntry.Message)
		s.ChannelMessageSend(m.ChannelID, message)
	}

}
