package experience_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/Luthor91/Tenshi/utils"
	"github.com/bwmarrin/discordgo"
)

// XPCommand gère les opérations d'expérience pour les utilisateurs
func ExperienceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%sexperience", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%sxp", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	// Extraire les arguments de la commande
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m]")
		return
	}

	// Affiche l'aide si -h est spécifié
	for _, arg := range parsedArgs {
		if arg.Arg == "-h" {
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m]")
			return
		}
	}

	// Si aucun argument n'est fourni, affiche l'XP de l'utilisateur qui a exécuté la commande
	if len(parsedArgs) == 0 {
		handleGetXP(s, m, m.Author.ID)
		return
	}

	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)

	// Variables pour stocker les informations sur l'utilisateur cible et l'action
	var targetUserID string
	var xpAmount int
	var action string

	// Analyser les arguments extraits
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			targetUserID = discord.HandleTarget(s, m, arg.Value).ID
		case "-r", "-s", "-a":
			if xpAmountValue, err := strconv.Atoi(arg.Value); err == nil {
				xpAmount = xpAmountValue
			} else {
				s.ChannelMessageSend(m.ChannelID, "Quantité d'XP invalide.")
				return
			}
			if arg.Arg == "-r" {
				if isMod || targetUserID == m.Author.ID {
					action = "remove"
				} else {
					s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour retirer de l'XP.")
					return
				}
			} else if arg.Arg == "-s" {
				if isMod {
					action = "set"
				} else {
					s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission pour définir l'XP.")
					return
				}
			} else if arg.Arg == "-a" {
				action = "add"
			}
		case "-m":
			action = "me"
		case "-g":
			if targetUserID != "" {
				action = "give"
			} else {
				s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un utilisateur à qui donner de l'XP.")
			}
		}
	}

	// Gérer les actions selon l'argument
	switch action {
	case "remove":
		handleRemoveXP(s, m, targetUserID, xpAmount)
	case "set":
		handleSetXP(s, m, targetUserID, xpAmount)
	case "add":
		handleAddXP(s, m, targetUserID, xpAmount)
	case "me":
		handleGetXP(s, m, m.Author.ID) // Pour obtenir l'XP de l'utilisateur qui a exécuté la commande
	case "give":
		handleGiveXP(s, m, m.Author.ID, targetUserID, xpAmount) // Giver ID is the command author
	default:
		s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, -m, ou -h.")
	}
}

// handleRemoveXP retire une quantité d'XP à l'utilisateur spécifié
func handleRemoveXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	userService := services.NewUserService()
	user, err := userService.GetUserByDiscordID(userID)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Utilisateur introuvable")
		return
	}

	err = userService.AddExperience(user, -amount)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Impossible de réduire l'XP")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s réduite de %d.", user.Username, amount))
}

func handleSetXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	userService := services.NewUserService()
	err := userService.SetExperience(userID, amount)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la définition de l'XP")
		return
	}
	user, _ := s.User(userID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s définie à %d.", user.Username, amount))
}

func handleAddXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string, amount int) {
	userService := services.NewUserService()
	user, err := userService.GetUserByDiscordID(userID)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Utilisateur introuvable")
		return
	}

	err = userService.AddExperience(user, amount)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Impossible d'ajouter de l'XP")
		return
	}

	discordUser, _ := s.User(userID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous avez gagné %d expérience.", discordUser.Username, amount))
}

// handleGetXP affiche l'XP de l'utilisateur spécifié
func handleGetXP(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	userService := services.NewUserService()
	amount, err := userService.GetExperience(userID)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Utilisateur non trouvé ou erreur lors de la récupération de l'XP")
		return
	}
	user, _ := s.User(userID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre expérience est de %d.", user.Username, amount))
}

// handleGiveXP permet à un utilisateur de donner une quantité d'XP à un autre utilisateur
func handleGiveXP(s *discordgo.Session, m *discordgo.MessageCreate, giverID, receiverID string, amount int) {
	if amount <= 0 {
		utils.SendErrorMessage(s, m.ChannelID, "La quantité d'XP donnée doit être supérieure à zéro")
		return
	}

	userService := services.NewUserService()
	err := userService.GiveMoney(giverID, receiverID, amount)
	if err != nil {
		utils.SendErrorMessage(s, m.ChannelID, "Pas assez d'XP pour faire le don ou utilisateur introuvable")
		return
	}

	giver, _ := s.User(giverID)
	receiver, _ := s.User(receiverID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s a donné %d XP à %s.", giver.Username, amount, receiver.Username))
}
