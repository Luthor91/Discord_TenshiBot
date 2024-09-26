package affinity_commands

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

// Structure temporaire pour stocker les modifications d'affinité
var tempAffinityMap = make(map[string]int) // Clé: userID, Valeur: changement temporaire d'affinité

// AffinityCommand gère les opérations d'affinité pour les utilisateurs
func AffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Commandes valides
	command := fmt.Sprintf("%saffinity", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%saff", config.AppConfig.BotPrefix)

	if !strings.HasPrefix(m.Content, command) && !strings.HasPrefix(m.Content, commandAlias) {
		return
	}

	// Extraire les arguments en utilisant ExtractArguments
	parsedArgs, err := discord.ExtractArguments(m.Content, command)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?affinity [-n <utilisateur>] [-r <quantité>] [-s <quantité>] [-a <quantité>] [-g] [-t] [-v]")
		return
	}

	// Si aucun argument n'est passé, afficher l'affinity de l'utilisateur qui a envoyé la commande
	if len(parsedArgs) == 0 {
		userService := services.NewUserService() // Crée une instance du UserService
		amount, err := userService.GetAffinity(m.Author.ID)
		if err != nil {
			utils.SendErrorMessage(s, m.ChannelID, "Erreur lors de la récupération de l'XP.")
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre affinité est de %d.", m.Author.Username, amount))
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)

	// Variables pour stocker les informations sur l'utilisateur cible, l'action et les options
	var targetUserID string
	var affinityAmount int
	var action string
	var verbose, temporary bool

	// Parcourir les arguments et les traiter
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			targetUserID = discord.HandleTarget(s, m, arg.Value).ID
		case "-r", "-s", "-a":
			affinityAmount, err = strconv.Atoi(arg.Value)
			if err != nil || affinityAmount < 0 {
				if verbose {
					s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'affinité.")
				}
				return
			}
			if arg.Arg == "-r" {
				if isMod || targetUserID == m.Author.ID {
					action = "remove"
				} else {
					if verbose {
						s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission de retirer de l'affinité.")
					}
					return
				}
			} else if arg.Arg == "-s" {
				if isMod {
					action = "set"
				} else {
					if verbose {
						s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission de définir l'affinité.")
					}
					return
				}
			} else if arg.Arg == "-a" {
				action = "add"
			}
		case "-g":
			action = "get"
		case "-v":
			verbose = true
		case "-t":
			temporary = true
		default:
			if verbose {
				s.ChannelMessageSend(m.ChannelID, "Argument non valide. Utilisez -n, -r, -s, -a, -g, -t, ou -v.")
			}
			return
		}
	}

	// Déterminer l'utilisateur cible (si aucun, utiliser l'utilisateur actuel)
	if targetUserID == "" {
		targetUserID = m.Author.ID
	}

	// Effectuer l'action appropriée
	switch action {
	case "remove":
		if temporary {
			handleTempAffinityChange(targetUserID, -affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité temporaire de %s réduite de %d.", targetUserID, affinityAmount))
			}
		} else {
			affinityService := services.NewAffinityService()
			_ = affinityService.AddAffinity(targetUserID, -affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s réduite de %d.", targetUserID, affinityAmount))
			}
		}
	case "set":
		if temporary {
			handleTempAffinitySet(targetUserID, affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité temporaire de %s définie à %d.", targetUserID, affinityAmount))
			}
		} else {
			affinityService := services.NewAffinityService()
			_ = affinityService.SetAffinity(targetUserID, affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s définie à %d.", targetUserID, affinityAmount))
			}
		}
	case "add":
		if temporary {
			handleTempAffinityChange(targetUserID, affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité temporaire de %s augmentée de %d.", targetUserID, affinityAmount))
			}
		} else {
			affinityService := services.NewAffinityService()
			_ = affinityService.AddAffinity(targetUserID, affinityAmount)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité de %s augmentée de %d.", targetUserID, affinityAmount))
			}
		}
	case "get":
		if temporary {
			tempAffinity := tempAffinityMap[targetUserID]
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité temporaire actuelle de %s : %d.", targetUserID, tempAffinity))
			}
		} else {
			affinityService := services.NewAffinityService()
			affinity, _ := affinityService.GetAffinity(targetUserID)
			if verbose {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Affinité actuelle de %s : %d.", targetUserID, affinity))
			}
		}
	default:
		if verbose {
			s.ChannelMessageSend(m.ChannelID, "Aucune action spécifiée. Utilisez -r, -s, -a, ou -g.")
		}
	}
}

// handleTempAffinityChange gère les changements temporaires d'affinité
func handleTempAffinityChange(userID string, amount int) {
	tempAffinityMap[userID] += amount
	go resetTemporaryAffinity(userID, amount, 5*time.Minute) // Réinitialisation après 5 minutes (ou tout autre délai souhaité)
}

// handleTempAffinitySet gère la définition temporaire d'affinité
func handleTempAffinitySet(userID string, amount int) {
	tempAffinityMap[userID] = amount
	go resetTemporaryAffinity(userID, amount, 5*time.Minute) // Réinitialisation après 5 minutes
}

// resetTemporaryAffinity réinitialise les changements d'affinité après une période donnée
func resetTemporaryAffinity(userID string, amount int, duration time.Duration) {
	time.Sleep(duration)
	tempAffinityMap[userID] -= amount
}
