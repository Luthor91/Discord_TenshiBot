package experience_commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m] [-t <durée en secondes>] [-v]")
		return
	}

	// Si aucun argument n'est passé, afficher l'XP de l'utilisateur qui a envoyé la commande
	if len(parsedArgs) == 0 {
		userService := services.NewUserService()
		amount, err := userService.GetExperience(m.Author.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre expérience est de %d.", m.Author.Username, amount))
		return
	}

	// Variables pour les actions et options
	var targetUser *discordgo.User
	var xpAmount int
	var action string
	var verbose bool
	var duration time.Duration

	// Affiche l'aide si -h est spécifié
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-h":
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m] [-t <durée en secondes>] [-v]")
			return
		case "-v":
			verbose = true
		case "-n":
			targetUser = discord.HandleTarget(s, m, arg.Value)
		case "-r", "-s", "-a":
			if xpAmountValue, err := strconv.Atoi(arg.Value); err == nil {
				xpAmount = xpAmountValue
				fmt.Println("XP Amount:", xpAmount) // Ligne de debug
			} else {
				if verbose {
					s.ChannelMessageSend(m.ChannelID, "Quantité d'XP invalide.")
				}
				return
			}
			if arg.Arg == "-r" {
				action = "remove"
			} else if arg.Arg == "-s" {
				action = "set"
			} else if arg.Arg == "-a" {
				action = "add"
			}
		case "-t":
			duration = time.Duration(arg.Duration.Seconds())
			return

		case "-g":
			action = "give"
		}
	}

	// Déterminer l'utilisateur cible (si aucun, utiliser l'utilisateur actuel)
	if targetUser == nil {
		targetUser = m.Author
	}

	// Gérer les actions selon l'argument
	userService := services.NewUserService()
	switch action {
	case "remove":
		originalXP, err := userService.GetExperience(targetUser.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP actuelle.")
			return
		}

		err = userService.AddExperience(targetUser.ID, -xpAmount)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Impossible de réduire l'XP")
			return
		}

		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s réduite de %d.", targetUser.Username, xpAmount))
		}

		if duration > 0 {
			go resetXPAfterDuration(targetUser.ID, originalXP, duration, s, m.ChannelID, verbose)
		}
	case "set":
		originalXP, err := userService.GetExperience(targetUser.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP actuelle.")
			return
		}

		err = userService.SetExperience(targetUser.ID, xpAmount)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la définition de l'XP")
			return
		}

		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s définie à %d.", targetUser.Username, xpAmount))
		}

		if duration > 0 {
			go resetXPAfterDuration(targetUser.ID, originalXP, duration, s, m.ChannelID, verbose)
		}
	case "add":
		fmt.Println("Action: add") // Ligne de debug
		originalXP, err := userService.GetExperience(targetUser.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP actuelle.")
			return
		}

		err = userService.AddExperience(targetUser.ID, xpAmount)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Impossible d'ajouter de l'XP")
			return
		}

		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s augmentée de %d.", targetUser.Username, xpAmount))
		}

		if duration > 0 {
			go resetXPAfterDuration(targetUser.ID, originalXP, duration, s, m.ChannelID, verbose)
		}
	case "give":
		originalXP, err := userService.GetExperience(targetUser.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP actuelle.")
			return
		}

		err = userService.GiveExperience(m.Author.ID, targetUser.ID, xpAmount)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Pas assez d'XP pour faire le don ou utilisateur introuvable")
			return
		}
		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s a donné %d XP à %s.", m.Author.Username, xpAmount, targetUser.Username))
		}

		if duration > 0 {
			user, _ := services.NewUserService().GetUserByDiscordID(m.Author.ID)
			go resetXPAfterDuration(user.UserDiscordID, user.Experience+xpAmount, duration, s, m.ChannelID, verbose)
			go resetXPAfterDuration(targetUser.ID, originalXP, duration, s, m.ChannelID, verbose)
		}

	default:
		if verbose {
			s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, -g, ou -h.")
		}
	}
}

// resetXPAfterDuration réinitialise l'XP après une durée donnée
func resetXPAfterDuration(userID string, originalXP int, duration time.Duration, s *discordgo.Session, channelID string, verbose bool) {
	time.Sleep(duration)
	userService := services.NewUserService()
	err := userService.SetExperience(userID, originalXP)
	if err != nil {
		utils.SendErrorMessage(s, channelID, "Erreur lors de la réinitialisation de l'XP.")
		return
	}

	if verbose {
		s.ChannelMessageSend(channelID, fmt.Sprintf("L'XP de l'utilisateur %s a été réinitialisée à %d après %v.", userID, originalXP, duration))
	}
}
