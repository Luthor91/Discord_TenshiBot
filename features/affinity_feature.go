package features

import (
	"strings"

	"github.com/Luthor91/Tenshi/models"
	"github.com/bwmarrin/discordgo"
)

// AdjustAffinity ajuste l'affinité d'un utilisateur en fonction du contenu de son message
func AdjustAffinity(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Empêcher le bot de répondre à lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est déjà dans la base
	user, exists := usersMap[m.Author.ID]
	if !exists {
		user = models.User{
			Username: m.Author.Username,
			Affinity: 0, // Affinité de départ
		}
	}

	// Vérifier les mots interdits
	for _, word := range badwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity-- // Diminuer l'affinité
			usersMap[m.Author.ID] = user
			SaveUsers()
			return
		}
	}

	// Vérifier les bons mots
	for _, word := range goodwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity++ // Augmenter l'affinité
			usersMap[m.Author.ID] = user
			SaveUsers()
			return
		}
	}
}

// Récupérer l'affinité d'un utilisateur
func GetUserAffinity(userID string) (models.User, bool) {
	LoadUsers() // Charger ou recharger les utilisateurs

	user, exists := usersMap[userID]
	return user, exists
}
