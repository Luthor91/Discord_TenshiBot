package controllers

import (
	"time"

	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// WarnController gère les opérations CRUD pour les Warns
type WarnController struct {
	DB *gorm.DB
}

// NewWarnController crée une nouvelle instance de WarnController
func NewWarnController() *WarnController {
	return &WarnController{DB: database.DB}
}

// CreateWarn crée un nouveau warn pour un utilisateur
func (wc *WarnController) CreateWarn(userDiscordID, reason, adminID string) error {
	warn := models.Warn{
		UserDiscordID: userDiscordID,
		Reason:        reason,
		WarnedAt:      time.Now(),
		AdminID:       adminID,
	}
	return wc.DB.Create(&warn).Error
}

// GetWarnsByUser récupère tous les warns d'un utilisateur
func (wc *WarnController) GetWarnsByUserDiscordID(userDiscordID string) ([]models.Warn, error) {
	var warns []models.Warn
	err := wc.DB.Where("user_discord_id = ?", userDiscordID).Find(&warns).Error
	return warns, err
}

// DeleteWarn supprime un warn par son ID
func (wc *WarnController) DeleteWarn(warnID uint) error {
	return wc.DB.Delete(&models.Warn{}, warnID).Error
}
