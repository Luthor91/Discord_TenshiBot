package controllers

import (
	"github.com/Luthor91/Tenshi/models"
	"gorm.io/gorm"
)

// WordController est un contrôleur pour gérer les mots positifs et négatifs
type WordController struct {
	DB *gorm.DB
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
	goodWord := models.GoodWord{
		Word: word,
	}
	// Vérifie si le mot existe déjà
	var existingWord models.GoodWord
	if err := ctrl.DB.Where("word = ?", word).First(&existingWord).Error; err == nil {
		return nil // Le mot existe déjà, rien à faire
	} else if err != gorm.ErrRecordNotFound {
		return err // Erreur inattendue
	}
	// Ajoute le nouveau mot
	return ctrl.DB.Create(&goodWord).Error
}

// AddBadWord ajoute un nouveau mauvais mot
func (ctrl *WordController) AddBadWord(word string) error {
	badWord := models.BadWord{
		Word: word,
	}
	// Vérifie si le mot existe déjà
	var existingWord models.BadWord
	if err := ctrl.DB.Where("word = ?", word).First(&existingWord).Error; err == nil {
		return nil // Le mot existe déjà, rien à faire
	} else if err != gorm.ErrRecordNotFound {
		return err // Erreur inattendue
	}
	// Ajoute le nouveau mot
	return ctrl.DB.Create(&badWord).Error
}
