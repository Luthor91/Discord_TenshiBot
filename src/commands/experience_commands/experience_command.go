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

// Structure temporaire pour stocker les modifications d'XP
var tempXPMap = make(map[string]int) // Clé: userID, Valeur: changement temporaire d'XP

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
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m] [-t] [-v]")
		return
	}

	// Variables pour les actions et options
	var targetUserID string
	var xpAmount int
	var action string
	var verbose, temporary bool

	// Affiche l'aide si -h est spécifié
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-h":
			s.ChannelMessageSend(m.ChannelID, "Utilisation : ?xp [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-m] [-t] [-v]")
			return
		case "-v":
			verbose = true
		case "-t":
			temporary = true
		case "-n":
			targetUserID = discord.HandleTarget(s, m, arg.Value).ID
		case "-r", "-s", "-a":
			if xpAmountValue, err := strconv.Atoi(arg.Value); err == nil {
				xpAmount = xpAmountValue
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
		case "-g":
			action = "give"
		case "-m":
			action = "me"
		}
	}

	// Déterminer l'utilisateur cible (si aucun, utiliser l'utilisateur actuel)
	if targetUserID == "" {
		targetUserID = m.Author.ID
	}

	// Gérer les actions selon l'argument
	userService := services.NewUserService() // Crée une instance du UserService
	switch action {
	case "remove":
		if temporary {
			handleTempXPChange(targetUserID, -xpAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP temporaire de %s réduite de %d.", targetUserID, xpAmount))
			}
		} else {
			err = userService.AddExperience(targetUserID, -xpAmount)
			if err != nil {
				utils.SendErrorMessage(s, m.ChannelID, "Impossible de réduire l'XP")
				return
			}

			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s réduite de %d.", targetUserID, xpAmount))
			}
		}
	case "set":
		if temporary {
			handleTempXPSet(targetUserID, xpAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP temporaire de %s définie à %d.", targetUserID, xpAmount))
			}
		} else {
			err = userService.SetExperience(targetUserID, xpAmount)
			if err != nil {
				utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la définition de l'XP")
				return
			}
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP de %s définie à %d.", targetUserID, xpAmount))
			}
		}
	case "add":
		if temporary {
			handleTempXPChange(targetUserID, xpAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("XP temporaire de %s augmentée de %d.", targetUserID, xpAmount))
			}
		} else {
			err = userService.AddExperience(targetUserID, xpAmount)
			if err != nil {
				utils.SendErrorMessage(s, m.ChannelID, "Impossible d'ajouter de l'XP")
				return
			}
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, vous avez gagné %d expérience.", targetUserID, xpAmount))
			}
		}
	case "me":
		amount, err := userService.GetExperience(m.Author.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Utilisateur non trouvé ou erreur lors de la récupération de l'XP")
			return
		}
		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre expérience est de %d.", m.Author.Username, amount))
		}
	case "give":
		err := userService.GiveExperience(m.Author.ID, targetUserID, xpAmount) // Correction de GiveMoney à GiveExperience
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Pas assez d'XP pour faire le don ou utilisateur introuvable")
			return
		}
		if verbose {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s a donné %d XP à %s.", m.Author.Username, xpAmount, targetUserID))
		}
	default:
		if verbose {
			s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, -m, ou -h.")
		}
	}
}

// handleTempXPChange gère les changements temporaires d'XP
func handleTempXPChange(userID string, changeAmount int) {
	tempXPMap[userID] += changeAmount
	go resetTemporaryXP(userID, 5*time.Minute) // Réinitialisation après 5 minutes (ou tout autre délai souhaité)
}

// handleTempXPSet gère la définition temporaire d'XP
func handleTempXPSet(userID string, setAmount int) {
	tempXPMap[userID] = setAmount
	go resetTemporaryXP(userID, 5*time.Minute) // Réinitialisation après 5 minutes
}

// resetTemporaryXP réinitialise les changements d'XP après une période donnée
func resetTemporaryXP(userID string, duration time.Duration) {
	time.Sleep(duration)
	tempXPMap[userID] = 0
}
