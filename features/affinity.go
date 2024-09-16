package features

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AdjustAffinity ajuste l'affinité d'un utilisateur en fonction du contenu de son message
func AdjustAffinity(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Empêcher le bot de répondre à lui-même
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est déjà dans la base
	user, exists := users[m.Author.ID]
	if !exists {
		user = User{
			Username: m.Author.Username,
			Affinity: 0, // Affinité de départ
		}
	}

	// Vérifier les mots interdits
	for _, word := range banwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity-- // Diminuer l'affinité
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s : %d (mot interdit utilisé)", user.Username, user.Affinity))
			users[m.Author.ID] = user
			SaveUsers()
			return
		}
	}

	// Vérifier les bons mots
	for _, word := range goodwords {
		if strings.Contains(strings.ToLower(m.Content), strings.ToLower(word)) {
			user.Affinity++ // Augmenter l'affinité
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s : %d (mot positif utilisé)", user.Username, user.Affinity))
			users[m.Author.ID] = user
			SaveUsers()
			return
		}
	}
}

// Récupérer l'affinité d'un utilisateur
func GetUserAffinity(userID string) (User, bool) {
	LoadUsers() // Charger ou recharger les utilisateurs

	user, exists := users[userID]
	return user, exists
}
