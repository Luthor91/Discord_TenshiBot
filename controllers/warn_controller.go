package controllers

import (
	"fmt"
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

// Compter le nombre de warns pour un utilisateur donné
func (wc *WarnController) CountWarnsByUser(userDiscordID string) (int64, error) {
	var count int64

	// Compter le nombre de warns dans la base de données pour cet utilisateur
	err := wc.DB.Where("user_discord_id = ?", userDiscordID).Count(&count).Error
	return count, err
}

// DeleteAllWarnsByUser supprime tous les avertissements pour un utilisateur donné
func (wc *WarnController) ResetWarns(userDiscordID string) error {
	result := wc.DB.Where("user_discord_id = ?", userDiscordID).Delete(&models.Warn{})
	if result.Error != nil {
		return fmt.Errorf("erreur lors de la suppression des avertissements: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("aucun avertissement trouvé pour l'utilisateur %s", userDiscordID)
	}
	return nil
}

// DeleteWarn supprime un warn par son ID
func (wc *WarnController) DeleteWarn(warnID uint) error {
	return wc.DB.Delete(&models.Warn{}, warnID).Error
}
