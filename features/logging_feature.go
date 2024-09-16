package features

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

var logMutex sync.Mutex

// LogMessage enregistre un message dans logs.json
func LogMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Récupérer les informations sur le serveur et le canal
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des informations du canal: %v", err)
		return
	}

	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Printf("Erreur lors de la récupération des informations du serveur: %v", err)
		return
	}

	// Créer une nouvelle entrée de log
	entry := models.LogEntry{
		Timestamp:   time.Now(),
		ServerID:    m.GuildID,
		ServerName:  guild.Name,
		ChannelID:   m.ChannelID,
		ChannelName: channel.Name,
		UserID:      m.Author.ID,
		Username:    m.Author.Username,
		Message:     m.Content,
	}

	// Lire les logs existants depuis logs.json
	logFile := "resources/logs.json"
	var logs []models.LogEntry
	data, err := os.ReadFile(logFile)
	if err == nil {
		json.Unmarshal(data, &logs)
	}

	// Ajouter le nouveau log
	logs = append(logs, entry)

	// Écrire les logs mis à jour dans logs.json
	data, err = json.MarshalIndent(logs, "", "  ")
	if err != nil {
		log.Printf("Erreur lors de la sérialisation des logs: %v", err)
		return
	}

	err = os.WriteFile(logFile, data, 0644)
	if err != nil {
		log.Printf("Erreur lors de l'écriture des logs dans le fichier: %v", err)
	}
}
