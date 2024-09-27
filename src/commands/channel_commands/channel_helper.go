package channel_commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/DiscordBot/services"
	"github.com/Luthor91/DiscordBot/utils"
	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles si seul ?channel est utilisé
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := `
**Commande ?channel :**

Arguments disponibles :
- **-n [nom]** : Spécifier le nom du salon (par défaut "channel").
- **-v** : Spécifier que c'est un salon vocal.
- **-t [durée]** : Spécifier la durée avant la suppression du salon (exemple : 1h, 30m).
- **-l** : Verrouiller ou déverrouiller le salon.
- **-c** : Créer un nouveau salon.
- **-d** : Supprimer un salon.
`
	s.ChannelMessageSend(channelID, helpMessage)
}

// Récupérer et analyser les arguments de la commande
func parseChannelArgs(m *discordgo.MessageCreate) (string, time.Duration, bool, bool, bool, bool, int, error) {
	args := strings.Fields(m.Content)
	var (
		duration      time.Duration
		channelName   string
		isVoice       bool
		shouldLock    bool
		createChannel bool
		deleteChannel bool
		messageCount  int // Nouvelle variable pour le nombre de messages à récupérer
		err           error
	)

	for i, arg := range args {
		switch arg {
		case "-t":
			if i+1 < len(args) {
				duration, err = utils.ParseDuration(args[i+1])
				if err != nil {
					return "", 0, false, false, false, false, 0, fmt.Errorf("temps non valide")
				}
			}
		case "-n":
			if i+1 < len(args) {
				channelName = args[i+1]
			}
		case "-v":
			isVoice = true
		case "-l":
			shouldLock = true
		case "-c":
			createChannel = true
		case "-d":
			deleteChannel = true
		case "-a":
			if i+1 < len(args) {
				messageCount, err = strconv.Atoi(args[i+1]) // Convertir le nombre de messages en entier
				if err != nil {
					return "", 0, false, false, false, false, 0, fmt.Errorf("nombre de messages non valide")
				}
			}
		}
	}

	if channelName == "" {
		channelName = "channel"
	}

	return channelName, duration, isVoice, shouldLock, createChannel, deleteChannel, messageCount, nil
}

// archiveMessages récupère les derniers messages d'un salon et les archive dans la base de données
func archiveMessages(s *discordgo.Session, m *discordgo.MessageCreate, archiveMessagesCount int) error {
	// Vérifie que le nombre de messages à archiver est valide
	if archiveMessagesCount <= 0 {
		return fmt.Errorf("le nombre de messages à archiver doit être supérieur à zéro")
	}

	// Récupérer les derniers messages du salon
	messages, err := s.ChannelMessages(m.ChannelID, archiveMessagesCount, "", "", "")
	if err != nil {
		return fmt.Errorf("erreur lors de la récupération des messages : %w", err)
	}

	// Créer une instance du service de log
	logService := services.NewLogService()

	// Archive chaque message en utilisant le service de log
	for _, msg := range messages {
		err = logService.InsertLog(s, msg) // Enregistrer le message dans la base de données
		if err != nil {
			return fmt.Errorf("erreur lors de l'enregistrement du message dans la base de données : %w", err)
		}
	}

	return nil
}
