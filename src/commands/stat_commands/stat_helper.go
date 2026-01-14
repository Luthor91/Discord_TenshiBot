package stat_commands

import (
	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles si seul ?channel est utilisé
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := `
**Commande ?stats :**

Arguments disponibles :
- **-u** : Affiche les statistiques de l'utilisateur.
- **-s** : Affiche les statistiques du serveur.
- **-b** : Affiche les statistiques du bot.
- **-c** : Affiche les statistiques du canal.
- **-h** : Affiche l'aide.

Exemple :
- ?stats -u
`
	s.ChannelMessageSend(channelID, helpMessage)
}
