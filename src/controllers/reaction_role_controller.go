package controllers

import (
	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/models"
	"gorm.io/gorm"
)

type ReactionRoleController struct {
	DB *gorm.DB
}

func NewReactionRoleController() *ReactionRoleController {
	return &ReactionRoleController{DB: database.DB}
}

func (rc *ReactionRoleController) Create(rr *models.ReactionRole) error {
	return rc.DB.Create(rr).Error
}

func (rc *ReactionRoleController) GetByMessageID(messageID string) ([]models.ReactionRole, error) {
	var roles []models.ReactionRole
	err := rc.DB.Where("message_id = ?", messageID).Find(&roles).Error
	return roles, err
}

func (rc *ReactionRoleController) DeleteByMessageID(messageID string) error {
	return rc.DB.Where("message_id = ?", messageID).Delete(&models.ReactionRole{}).Error
}
