package moderation_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/bwmarrin/discordgo"
)

// ModerateMessageCommand gère les commandes de modération des messages
func ModerateMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := fmt.Sprintf("%smessage", config.AppConfig.BotPrefix)
	commandAlias := fmt.Sprintf("%smsg", config.AppConfig.BotPrefix)

	var prefix string
	if strings.HasPrefix(m.Content, command) {
		prefix = command
	} else if strings.HasPrefix(m.Content, commandAlias) {
		prefix = commandAlias
	} else {
		return
	}

	// Extraction des arguments
	parsedArgs, err := discord.ExtractArguments(m.Content, prefix)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	// Initialiser avec le salon actuel comme valeur par défaut
	var userID string = ""
	var channelID string = m.ChannelID // Utiliser le salon actuel par défaut
	var deleteCount int
	var verbose bool

	// Analyse des arguments extraits
	for _, arg := range parsedArgs {
		switch arg.Arg {
		case "-n":
			// Résoudre l'ID utilisateur à partir du pseudo Discord
			if arg.Value != "" {
				userID, err = resolveUserID(s, m.GuildID, arg.Value)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la résolution du pseudo : %v", err))
					return
				}
			}
		case "-c":
			// Si l'argument -c est présent mais vide, on garde le salon actuel
			if arg.Value != "" {
				// Résoudre l'ID du salon à partir du nom ou ID
				tempChannelID, err := resolveChannelID(s, m.GuildID, arg.Value)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la résolution du salon : %v", err))
					return
				}
				channelID = tempChannelID
			}
		case "-d":
			// Vérifier et convertir la valeur de suppression
			deleteCount, err = strconv.Atoi(arg.Value)
			if err != nil || deleteCount <= 0 {
				s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre valide de messages à supprimer.")
				return
			}
		case "-v":
			verbose = true
		default:
			// Ignorer les arguments non reconnus
		}
	}

	// Validation de l'input
	if deleteCount == 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez spécifier un nombre de messages à supprimer avec -d.")
		return
	}

	// Debug - Afficher les valeurs des arguments
	fmt.Printf("Suppression de %d messages dans le salon ID: %s\n", deleteCount, channelID)

	// Récupérer les messages à supprimer
	messages, err := s.ChannelMessages(channelID, deleteCount, "", "", "")
	if err != nil {
		errMsg := fmt.Sprintf("Erreur lors de la récupération des messages: %v", err)
		fmt.Println(errMsg)
		s.ChannelMessageSend(m.ChannelID, errMsg)
		return
	}

	fmt.Printf("Nombre de messages récupérés: %d\n", len(messages))

	// Supprimer les messages selon l'utilisateur ou le canal spécifié
	var deletedCount int = 0

	// Pour la suppression en masse (bulk delete)
	var messageIDs []string

	for _, message := range messages {
		if userID == "" || message.Author.ID == userID {
			messageIDs = append(messageIDs, message.ID)
			deletedCount++
		}
	}

	// Si nous avons des messages à supprimer
	if len(messageIDs) > 0 {
		// Pour les messages de moins de 2 semaines, on peut utiliser BulkDelete
		if len(messageIDs) > 1 {
			err = s.ChannelMessagesBulkDelete(channelID, messageIDs)
			if err != nil {
				fmt.Printf("Erreur lors de la suppression groupée des messages: %v\n", err)

				// Fallback: suppression individuelle si BulkDelete échoue
				for _, msgID := range messageIDs {
					err := s.ChannelMessageDelete(channelID, msgID)
					if err != nil {
						fmt.Printf("Erreur lors de la suppression du message %s: %v\n", msgID, err)
					}
				}
			}
		} else if len(messageIDs) == 1 {
			// Un seul message à supprimer
			err := s.ChannelMessageDelete(channelID, messageIDs[0])
			if err != nil {
				fmt.Printf("Erreur lors de la suppression du message: %v\n", err)
			}
		}
	}

	// Afficher le message de confirmation si l'option "-v" est présente
	if verbose {
		channel, err := s.Channel(channelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d messages supprimés, mais erreur lors de la récupération des informations du salon.", deletedCount))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d messages supprimés avec succès dans le salon %s.", deletedCount, channel.Name))
		}
	}
}

// resolveUserID résout un pseudo Discord en ID
func resolveUserID(s *discordgo.Session, guildID, username string) (string, error) {
	// Recherche de l'utilisateur dans le serveur
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return "", fmt.Errorf("impossible de récupérer les membres du serveur : %v", err)
	}

	// Recherche de l'utilisateur dans la liste des membres
	for _, member := range members {
		if strings.EqualFold(member.User.Username, username) {
			return member.User.ID, nil
		}
	}

	// Si l'utilisateur n'est pas trouvé, on retourne une erreur
	return "", fmt.Errorf("utilisateur non trouvé : %s", username)
}

// resolveChannelID résout un nom de salon ou un ID en ID de salon
func resolveChannelID(s *discordgo.Session, guildID, channelArg string) (string, error) {
	// Si l'argument commence par <# et se termine par >, c'est un ID de salon sous forme de mention
	if strings.HasPrefix(channelArg, "<#") && strings.HasSuffix(channelArg, ">") {
		return channelArg[2 : len(channelArg)-1], nil
	}

	// Vérifie si l'argument est directement un ID numérique
	if _, err := strconv.ParseUint(channelArg, 10, 64); err == nil {
		// Vérifie que l'ID existe bien sur le serveur
		channel, err := s.Channel(channelArg)
		if err == nil && channel.GuildID == guildID {
			return channelArg, nil
		}
	}

	// Si l'argument est un nom de salon, rechercher l'ID du salon
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return "", fmt.Errorf("impossible de récupérer les salons du serveur : %v", err)
	}

	// Extrait seulement les caractères alphanumériques du nom recherché
	// pour une correspondance basée uniquement sur le texte visible
	var alphaNumOnly = func(input string) string {
		var result strings.Builder
		for _, r := range input {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
				result.WriteRune(r)
			}
		}
		return result.String()
	}

	cleanedArg := alphaNumOnly(strings.ToLower(channelArg))
	fmt.Printf("Argument nettoyé: [%s] (len=%d)\n", cleanedArg, len(cleanedArg))

	// Deux passes: correspondance exacte, puis correspondance partielle
	for _, channel := range channels {
		// Vérification par texte visible uniquement (en ignorant les caractères spéciaux)
		cleanChannelName := alphaNumOnly(strings.ToLower(channel.Name))

		// Correspondance exacte sur la partie alphanumérique
		if cleanChannelName == cleanedArg {
			fmt.Printf("Correspondance exacte trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	// Deuxième passe: correspondance partielle
	for _, channel := range channels {
		cleanChannelName := alphaNumOnly(strings.ToLower(channel.Name))

		// Correspondance partielle
		if strings.Contains(cleanChannelName, cleanedArg) ||
			strings.Contains(cleanedArg, cleanChannelName) {
			fmt.Printf("Correspondance partielle trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	// Dernière tentative: correspondance phonétique ou visuelle proche
	for _, channel := range channels {
		// Affiche les noms avec leurs codes Unicode pour le débogage
		fmt.Printf("Canal: [%s] -> alphaNum: [%s]\n",
			channel.Name,
			alphaNumOnly(channel.Name))

		// Si le canal contient "ticket" quelque part et que l'argument contient aussi "ticket"
		if strings.Contains(strings.ToLower(channel.Name), "ticket") &&
			strings.Contains(strings.ToLower(channelArg), "ticket") {
			fmt.Printf("Correspondance par mot-clé trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("salon non trouvé : %s", channelArg)
}
