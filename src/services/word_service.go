package services

import (
	"github.com/Luthor91/DiscordBot/controllers"
	"gorm.io/gorm"
)

// WordService gère les opérations liées aux mots
type WordService struct {
	controller *controllers.WordController
}

func NewWordService(db *gorm.DB) *WordService {
	return &WordService{
		controller: controllers.NewWordController(db),
	}
}

// AddBadWord ajoute un mauvais mot
func (s *WordService) AddBadWord(word string) error {
	return s.controller.AddBadWord(word)
}

// DeleteBadWord supprime un mauvais mot
func (s *WordService) DeleteBadWord(word string) error {
	return s.controller.DeleteBadWord(word)
}

// ListBadWords retourne la liste des mauvais mots
func (s *WordService) ListBadWords() ([]string, error) {
	return s.controller.GetBadWords()
}
