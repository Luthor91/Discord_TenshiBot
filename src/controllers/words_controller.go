package controllers

import (
	"github.com/Luthor91/DiscordBot/models"
	"gorm.io/gorm"
)

// WordController gère les opérations sur les mots positifs et négatifs
type WordController struct {
	DB *gorm.DB
}

// NewWordController crée une nouvelle instance de WordController
func NewWordController(db *gorm.DB) *WordController {
	if db == nil {
		panic("WordController initialized with nil DB")
	}
	return &WordController{DB: db}
}


// GetBadWords récupère tous les mauvais mots
func (ctrl *WordController) GetBadWords() ([]string, error) {
	var words []models.BadWord
	if err := ctrl.DB.Find(&words).Error; err != nil {
		return nil, err
	}

	var badWords []string
	for _, word := range words {
		badWords = append(badWords, word.Word)
	}
	return badWords, nil
}

// AddBadWord ajoute un nouveau mauvais mot
func (ctrl *WordController) AddBadWord(word string) error {
	badWord := models.BadWord{Word: word}
	return ctrl.DB.Create(&badWord).Error
}

// DeleteBadWord supprime un mauvais mot par son nom
func (ctrl *WordController) DeleteBadWord(word string) error {
	res := ctrl.DB.Where("word = ?", word).Delete(&models.BadWord{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

