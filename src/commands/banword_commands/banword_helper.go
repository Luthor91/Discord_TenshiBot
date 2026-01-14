package banword_commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles si seul ?channel est utilisé
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : gestion des mots bannis (`?banword`)**\n\n" +

		"**Usage général :**\n" +
		"`?banword [options] <mot>`\n\n" +

		"**Options disponibles :**\n" +
		"- `-a` : ajouter un mot à la liste des mots bannis\n" +
		"- `-d` : supprimer un mot de la liste des mots bannis\n" +
		"- `-l` : afficher la liste des mots bannis\n" +
		"- `-v` : mode verbeux (affiche les confirmations)\n" +
		"- `-h` : afficher cette aide\n\n" +

		"**Exemples :**\n" +
		"- `?banword -a mot-a-ajouter`\n" +
		"- `?banword -d mot-a-supprimer`\n" +
		"- `?banword -l`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}


// handleBadWord gère l'ajout d'un mauvais mot
func handleBadWord(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	word string,
	verbose bool,
) {
	err := wordService.AddBadWord(word)
	if err != nil && verbose {
		// Les erreurs doivent rester visibles
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Erreur lors de l'ajout du mauvais mot : %s", err),
		)
		return
	}

	if verbose {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Ajout du mauvais mot : %s", word),
		)
	}
}

// handleDeleteWord gère la suppression d'un mot
func handleDeleteWord(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	word string,
	verbose bool,
) {
	err := wordService.DeleteBadWord(word)
	if err != nil && verbose {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Erreur lors de la suppression du mauvais mot : %s", err),
		)
		return
	}

	if verbose {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("Suppression du mot : %s", word),
		)
	}
}


// handleListBadWords affiche les mauvais mots
func handleListBadWords(s *discordgo.Session, m *discordgo.MessageCreate) {
	words, err := wordService.ListBadWords()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des mauvais mots.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Voici la liste des banwords : %v", words))
}
