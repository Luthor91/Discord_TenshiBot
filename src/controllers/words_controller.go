package controllers

import (
	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/models"
	"gorm.io/gorm"
)

// WordController gère les opérations sur les mots positifs et négatifs
type WordController struct {
	DB *gorm.DB
}

// NewWordController crée une nouvelle instance de WordController
func NewWordController() *WordController {
	return &WordController{
		DB: database.DB,
	}
}

// GetGoodWords récupère tous les bons mots
func (ctrl *WordController) GetGoodWords() ([]string, error) {
	var words []models.GoodWord
	if err := ctrl.DB.Find(&words).Error; err != nil {
		return nil, err
	}
	var goodWords []string
	for _, word := range words {
		goodWords = append(goodWords, word.Word)
	}
	return goodWords, nil
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

// AddGoodWord ajoute un nouveau bon mot
func (ctrl *WordController) AddGoodWord(word string) error {
	goodWord := models.GoodWord{Word: word}
	// Vérifie si le mot existe déjà
	if err := ctrl.DB.Where("word = ?", word).First(&models.GoodWord{}).Error; err == nil {
		return nil // Le mot existe déjà, rien à faire
	} else if err != gorm.ErrRecordNotFound {
		return err // Erreur inattendue
	}
	// Ajoute le nouveau mot
	return ctrl.DB.Create(&goodWord).Error
}

// AddBadWord ajoute un nouveau mauvais mot
func (ctrl *WordController) AddBadWord(word string) error {
	badWord := models.BadWord{Word: word}
	// Vérifie si le mot existe déjà
	if err := ctrl.DB.Where("word = ?", word).First(&models.BadWord{}).Error; err == nil {
		return nil // Le mot existe déjà, rien à faire
	} else if err != gorm.ErrRecordNotFound {
		return err // Erreur inattendue
	}
	// Ajoute le nouveau mot
	return ctrl.DB.Create(&badWord).Error
}

// DeleteGoodWord supprime un bon mot par son nom
func (ctrl *WordController) DeleteGoodWord(word string) error {
	// Vérifie si le mot existe
	if err := ctrl.DB.Where("word = ?", word).Delete(&models.GoodWord{}).Error; err != nil {
		return err // Erreur lors de la suppression
	}
	return nil // Suppression réussie
}

// DeleteBadWord supprime un mauvais mot par son nom
func (ctrl *WordController) DeleteBadWord(word string) error {
	// Vérifie si le mot existe
	if err := ctrl.DB.Where("word = ?", word).Delete(&models.BadWord{}).Error; err != nil {
		return err // Erreur lors de la suppression
	}
	return nil // Suppression réussie
}
