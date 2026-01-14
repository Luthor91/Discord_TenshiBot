package channel_commands

import (
	"fmt"

	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles pour la commande ?channel
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : gestion des salons (`?channel`)**\n\n" +

		"**Usage général :**\n" +
		"`?channel [options] <nom>`\n\n" +

		"**Options disponibles :**\n" +
		"- `-c`           : créer un nouveau salon\n" +
		"- `-d`           : supprimer un salon existant\n" +
		"- `-l`           : verrouiller ou déverrouiller un salon\n" +
		"- `-v`           : créer un salon vocal\n" +
		"- `-t <durée>`   : durée avant suppression ou verrouillage (ex : `30s`, `10m`, `1h`)\n" +
		"- `-h`           : afficher cette aide\n\n" +

		"**Paramètre :**\n" +
		"- `<nom>` : nom du salon cible (obligatoire pour `-c` et `-d`)\n\n" +

		"**Comportement :**\n" +
		"- `-c` : crée un salon (texte par défaut, vocal si `-v`)\n" +
		"- `-d` : supprime le salon spécifié\n" +
		"- `-l` : verrouille ou déverrouille le salon spécifié\n" +
		"- `-t` : optionnelle, utilisée pour la suppression différée ou le verrouillage\n\n" +

		"**Exemples :**\n" +
		"- `?channel -c -t 30s salon-temporaire`\n" +
		"- `?channel -c -v salon-vocal`\n" +
		"- `?channel -l salon-prive`\n" +
		"- `?channel -d salon-a-supprimer`\n"

	s.ChannelMessageSend(channelID, helpMessage)
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
