package affinity_commands

import (
	"fmt"
	"strconv"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// parseAffinityArgs analyse les arguments de la commande et retourne l'utilisateur cible, la quantité d'affinité, et l'action
func parseAffinityArgs(args []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, int, string, error) {
	var targetUserID string
	var affinityAmount int
	var action string

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "-n":
			if i+1 < len(args) {
				targetUserID = args[i+1]
				i++ // Sauter le nom de l'utilisateur
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier un utilisateur avec -n <utilisateur>")
			}
		case "-r":
			action = "remove"
			if i+1 < len(args) {
				var err error
				affinityAmount, err = strconv.Atoi(args[i+1])
				if err != nil || affinityAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'affinité")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -r <quantité>")
			}
		case "-s":
			action = "set"
			if i+1 < len(args) {
				var err error
				affinityAmount, err = strconv.Atoi(args[i+1])
				if err != nil || affinityAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'affinité")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -s <quantité>")
			}
		case "-a":
			action = "add"
			if i+1 < len(args) {
				var err error
				affinityAmount, err = strconv.Atoi(args[i+1])
				if err != nil || affinityAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'affinité")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -a <quantité>")
			}
		case "-g":
			action = "get"
		}
	}

	// Déterminer l'utilisateur cible
	if targetUserID != "" {
		isAdmin, _ := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
		if !isAdmin {
			return "", 0, "", fmt.Errorf("vous n'avez pas la permission de modifier l'affinité d'un autre utilisateur")
		}
	}

	if targetUserID == "" {
		targetUserID = m.Author.ID
	}

	return targetUserID, affinityAmount, action, nil
}

// handleRemoveAffinity retire une quantité d'affinité à l'utilisateur spécifié
func handleRemoveAffinity(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	user, exists := services.NewAffinityService(s).GetUserAffinity(userID)
	if !exists || user.Affinity < amount {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé ou pas assez d'affinité")
		return
	}
	newAffinity := user.Affinity - amount
	err := services.NewUserService().SetAffinity(userID, newAffinity)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'affinité")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s réduite à %d.", user.Username, newAffinity))
}

// handleSetAffinity définit une quantité d'affinité pour l'utilisateur spécifié
func handleSetAffinity(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	err := services.NewUserService().SetAffinity(userID, amount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la définition de l'affinité")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s définie à %d.", userID, amount))
}

// handleAddAffinity ajoute une quantité d'affinité à l'utilisateur spécifié
func handleAddAffinity(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	user, exists := services.NewAffinityService(s).GetUserAffinity(userID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé")
		return
	}
	newAffinity := user.Affinity + amount
	err := services.NewUserService().SetAffinity(userID, newAffinity)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout d'affinité")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s augmentée à %d.", user.Username, newAffinity))
}

// handleGetAffinity affiche l'affinité de l'utilisateur spécifié
func handleGetAffinity(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	user, exists := services.NewAffinityService(s).GetUserAffinity(userID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité actuelle de %s : %d", user.Username, user.Affinity))
}
