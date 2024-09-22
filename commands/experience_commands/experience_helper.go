package experience_commands

import (
	"fmt"
	"strconv"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// parseXPArgs analyse les arguments de la commande et retourne l'utilisateur cible, la quantité d'XP, et l'action
func parseXPArgs(args []string, m *discordgo.MessageCreate, s *discordgo.Session) (string, int, string, error) {
	var targetUserID string
	var xpAmount int
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
				xpAmount, err = strconv.Atoi(args[i+1])
				if err != nil || xpAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'XP")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -r <quantité>")
			}
		case "-s":
			action = "set"
			if i+1 < len(args) {
				var err error
				xpAmount, err = strconv.Atoi(args[i+1])
				if err != nil || xpAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'XP")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -s <quantité>")
			}
		case "-a":
			action = "add"
			if i+1 < len(args) {
				var err error
				xpAmount, err = strconv.Atoi(args[i+1])
				if err != nil || xpAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'XP")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -a <quantité>")
			}
		case "-m":
			action = "me"
		case "-g":
			action = "give"
			if i+1 < len(args) {
				var err error
				xpAmount, err = strconv.Atoi(args[i+1])
				if err != nil || xpAmount < 0 {
					return "", 0, "", fmt.Errorf("veuillez entrer une quantité valide d'XP")
				}
				i++
			} else {
				return "", 0, "", fmt.Errorf("veuillez spécifier une quantité avec -d <quantité>")
			}
		}
	}

	// Déterminer l'utilisateur cible
	if targetUserID != "" {
		isAdmin, _ := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
		if !isAdmin {
			return "", 0, "", fmt.Errorf("vous n'avez pas la permission de modifier l'XP d'un autre utilisateur")
		}
	}

	if targetUserID == "" {
		targetUserID = m.Author.ID
	}

	return targetUserID, xpAmount, action, nil
}

// handleRemoveXP retire une quantité d'XP à l'utilisateur spécifié
func handleRemoveXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	current, exists := services.NewUserService(controllers.NewUserController()).GetExperience(userID)
	if !exists || current < amount {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé ou pas assez d'XP")
		return
	}
	newXP := current - amount
	err := services.NewUserService(controllers.NewUserController()).SetExperience(userID, newXP)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors du retrait d'XP")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP s réduite à %d.", newXP))
}

// handleSetXP définit une quantité d'XP pour l'utilisateur spécifié
func handleSetXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	err := services.NewUserService(controllers.NewUserController()).SetExperience(userID, amount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la définition de l'XP")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s définie à %d.", userID, amount))
}

// handleAddXP ajoute une quantité d'XP à l'utilisateur spécifié
func handleAddXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	current, exists := services.NewUserService(controllers.NewUserController()).GetExperience(userID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé")
		return
	}
	newXP := current + amount
	err := services.NewUserService(controllers.NewUserController()).SetExperience(userID, newXP)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'ajout d'XP")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP augmentée à %d.", newXP))
}

// handleGetXP affiche l'XP de l'utilisateur spécifié
func handleGetXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	amount, exists := services.NewUserService(controllers.NewUserController()).GetExperience(userID)
	if !exists {
		s.ChannelMessageSend(m.ChannelID, "Erreur : utilisateur non trouvé")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP actuelle %d", amount))
}

// handleGiveXP permet à un utilisateur de donner une quantité d'XP à un autre utilisateur
func handleGiveXP(s *discordgo.Session, m *discordgo.MessageCreate, giverID, receiverID string, amount int) {
	currentGiverXP, exists := services.NewUserService(controllers.NewUserController()).GetExperience(giverID)
	if !exists || currentGiverXP < amount {
		s.ChannelMessageSend(m.ChannelID, "Erreur : pas assez d'XP pour faire le don.")
		return
	}

	// Retirer l'XP de l'utilisateur donneur
	newGiverXP := currentGiverXP - amount
	err := services.NewUserService(controllers.NewUserController()).SetExperience(giverID, newGiverXP)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la mise à jour de l'XP du donneur.")
		return
	}

	// Ajouter l'XP à l'utilisateur receveur
	currentReceiverXP, _ := services.NewUserService(controllers.NewUserController()).GetExperience(receiverID)
	newReceiverXP := currentReceiverXP + amount
	err = services.NewUserService(controllers.NewUserController()).SetExperience(receiverID, newReceiverXP)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la mise à jour de l'XP du receveur.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s a donné %d XP à %s. XP restante : %d.", m.Author.Username, amount, receiverID, newGiverXP))
}
