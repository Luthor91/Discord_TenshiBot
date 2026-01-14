package kick_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?kick
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : expulsion d'un utilisateur (`?kick`)**\n\n" +

		"**Usage :**\n" +
		"`?kick <utilisateur> [raison]`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?kick @user`\n" +
		"- `?kick @user comportement inapproprié`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
