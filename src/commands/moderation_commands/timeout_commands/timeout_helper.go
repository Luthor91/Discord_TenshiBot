package timeout_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?timeout
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : mise en timeout d'un utilisateur (`?timeout`)**\n\n" +

		"**Usage :**\n" +
		"`?timeout <utilisateur> <durée> [raison]`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `durée`       : durée du timeout (ex: `10m`, `1h`, `1d`) (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?timeout @user 10m`\n" +
		"- `?timeout @user 1d comportement abusif`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
