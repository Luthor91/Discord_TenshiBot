package services

import (
	"log"
	"sync"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

// LogService est un service pour gérer les logs des messages
type LogService struct {
	logCtrl *controllers.LogController
	mu      sync.Mutex
}

// NewLogService crée une nouvelle instance de LogService
func NewLogService(logCtrl *controllers.LogController) *LogService {
	return &LogService{
		logCtrl: logCtrl,
	}
}

// LogMessage enregistre un message dans la base de données
func (service *LogService) LogMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	service.mu.Lock()
	defer service.mu.Unlock()

	// Récupérer les informations sur le serveur et le canal
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des informations du canal: %v", err)
		return err
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des informations du serveur: %v", err)
		return err
	}

	// Créer une nouvelle entrée de log
	entry := models.Log{
		Timestamp:     time.Now(),
		ServerID:      m.GuildID,
		ServerName:    guild.Name,
		ChannelID:     m.ChannelID,
		ChannelName:   channel.Name,
		UserDiscordID: m.Author.ID,
		Username:      m.Author.Username,
		Message:       m.Content,
	}

	// Enregistrer l'entrée de log dans la base de données
	err = service.logCtrl.SaveLog(entry)
	if err != nil {
		log.Printf("Erreur lors de l'enregistrement du log dans la base de données: %v", err)
		return err
	}

	return nil
}
