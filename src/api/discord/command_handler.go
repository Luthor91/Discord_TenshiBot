package discord

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Argument représente un argument sous la forme { arg : "-n", value : "value"}
type Argument struct {
	Arg      string
	Value    string
	Duration time.Duration // Si applicable, la durée pour l'argument "-t"
}

// Gérer le ciblage d'un utilisateur par son nom ou par mention
func HandleTarget(s *discordgo.Session, m *discordgo.MessageCreate, target string) *discordgo.User {
	// Vérifier d'abord les mentions dans le message
	for _, mention := range m.Mentions {
		if mention.Username == target || fmt.Sprintf("%s#%s", mention.Username, mention.Discriminator) == target {
			return mention
		}
	}

	// Ensuite, vérifier les utilisateurs dans le serveur
	users, err := s.GuildMembers(m.GuildID, "", 100)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des membres.")
		return nil
	}

	for _, user := range users {
		if user.User.Username == target || fmt.Sprintf("%s#%s", user.User.Username, user.User.Discriminator) == target {
			return user.User
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Utilisateur non trouvé.")
	return nil
}

// HandleChannel récupère le salon à partir d'une mention ou d'un nom donné.
func HandleChannel(s *discordgo.Session, m *discordgo.MessageCreate, target string) (*discordgo.Channel, error) {
	if strings.HasPrefix(target, "<#") && strings.HasSuffix(target, ">") {
		channelID := strings.TrimPrefix(strings.TrimSuffix(target, ">"), "<#")
		channel, err := s.Channel(channelID)
		if err != nil {
			return nil, fmt.Errorf("Salon mentionné introuvable.")
		}
		return channel, nil
	}
	// Si ce n'est pas une mention, on peut essayer de récupérer le salon par nom
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de la récupération des salons : %s", err.Error())
	}

	for _, channel := range channels {
		if channel.Name == target {
			return channel, nil
		}
	}

	return nil, fmt.Errorf("Salon avec le nom '%s' introuvable.", target)
}

// ExtractArguments récupère et valide les arguments d'un message
// et parse la durée spécifiée avec "-t".
func ExtractArguments(content, command string) ([]Argument, error) {
	args := strings.Fields(strings.TrimPrefix(content, command))
	if len(args) < 1 {
		return nil, nil
	}

	var parsedArgs []Argument
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			// Si l'argument suivant n'est pas une option (n'a pas de "-"), on le considère comme la valeur
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				if arg == "-t" {
					// Si l'argument est "-t", on parse la durée
					duration, err := parseDuration(args[i+1])
					if err != nil {
						return nil, err
					}
					parsedArgs = append(parsedArgs, Argument{Arg: arg, Value: args[i+1], Duration: duration})
				} else {
					// Sinon, on ajoute l'argument avec sa valeur
					parsedArgs = append(parsedArgs, Argument{Arg: arg, Value: args[i+1]})
				}
				i++ // On saute la valeur qui vient d'être utilisée
			} else {
				// L'argument n'a pas de valeur associée
				parsedArgs = append(parsedArgs, Argument{Arg: arg, Value: ""})
			}
		}
	}

	return parsedArgs, nil
}

// ParseDuration parse la durée au format '10s', '5m', '2h', '1d'
func parseDuration(durationStr string) (time.Duration, error) {
	if len(durationStr) < 2 {
		return 0, fmt.Errorf("durée invalide")
	}

	unit := durationStr[len(durationStr)-1]
	valueStr := durationStr[:len(durationStr)-1]
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("valeur numérique invalide")
	}

	switch unit {
	case 's': // secondes
		return time.Duration(value) * time.Second, nil
	case 'm': // minutes
		return time.Duration(value) * time.Minute, nil
	case 'h': // heures
		return time.Duration(value) * time.Hour, nil
	case 'd': // jours
		return time.Duration(value) * time.Hour * 24, nil
	default:
		return 0, fmt.Errorf("unité de temps invalide")
	}
}
