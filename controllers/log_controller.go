package controllers

import (
	"time"

	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// LogController est un contrôleur pour gérer les opérations sur les entrées de journal
type LogController struct {
	DB *gorm.DB
}

// CreateLog crée une nouvelle entrée de journal
func (ctrl *LogController) CreateLog(timestamp time.Time, serverID, serverName, channelID, channelName, userID, username, message string) (*models.Log, error) {
	logEntry := models.Log{
		Timestamp:     timestamp,
		ServerID:      serverID,
		ServerName:    serverName,
		ChannelID:     channelID,
		ChannelName:   channelName,
		UserDiscordID: userID,
		Username:      username,
		Message:       message,
	}
	if err := ctrl.DB.Create(&logEntry).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// GetLog récupère une entrée de journal par ID
func (ctrl *LogController) GetLog(id uint) (*models.Log, error) {
	var logEntry models.Log
	if err := ctrl.DB.First(&logEntry, id).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// UpdateLog met à jour une entrée de journal
func (ctrl *LogController) UpdateLog(id uint, timestamp time.Time, serverID, serverName, channelID, channelName, userID, username, message string) (*models.Log, error) {
	var logEntry models.Log
	if err := ctrl.DB.First(&logEntry, id).Error; err != nil {
		return nil, err
	}
	logEntry.Timestamp = timestamp
	logEntry.ServerID = serverID
	logEntry.ServerName = serverName
	logEntry.ChannelID = channelID
	logEntry.ChannelName = channelName
	logEntry.UserDiscordID = userID
	logEntry.Username = username
	logEntry.Message = message
	if err := ctrl.DB.Save(&logEntry).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}

// DeleteLog supprime une entrée de journal
func (ctrl *LogController) DeleteLog(id uint) error {
	if err := ctrl.DB.Delete(&models.Log{}, id).Error; err != nil {
		return err
	}
	return nil
}

// SaveLog enregistre une entrée de journal dans la base de données
func (ctrl *LogController) SaveLog(entry models.Log) error {
	if err := ctrl.DB.Create(&entry).Error; err != nil {
		return err
	}
	return nil
}
