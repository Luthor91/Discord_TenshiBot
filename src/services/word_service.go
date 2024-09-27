package services

import (
	"github.com/Luthor91/DiscordBot/controllers"
)

// WordService gère les opérations liées aux mots
type WordService struct {
	controller *controllers.WordController
}

// NewWordService crée une nouvelle instance de WordService
func NewWordService() *WordService {
	return &WordService{
		controller: controllers.NewWordController(),
	}
}

// AddGoodWord ajoute un bon mot
func (s *WordService) AddGoodWord(word string) error {
	return s.controller.AddGoodWord(word)
}

// AddBadWord ajoute un mauvais mot
func (s *WordService) AddBadWord(word string) error {
	return s.controller.AddBadWord(word)
}

// DeleteGoodWord supprime un bon mot
func (s *WordService) DeleteGoodWord(word string) error {
	return s.controller.DeleteGoodWord(word)
}

// DeleteBadWord supprime un mauvais mot
func (s *WordService) DeleteBadWord(word string) error {
	return s.controller.DeleteBadWord(word)
}

// ListGoodWords retourne la liste des bons mots
func (s *WordService) ListGoodWords() ([]string, error) {
	return s.controller.GetGoodWords()
}

// ListBadWords retourne la liste des mauvais mots
func (s *WordService) ListBadWords() ([]string, error) {
	return s.controller.GetBadWords()
}
