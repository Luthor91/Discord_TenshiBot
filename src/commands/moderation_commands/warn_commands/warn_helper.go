package warn_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?warn
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : avertissement d'un utilisateur (`?warn`)**\n\n" +

		"**Usage :**\n" +
		"`?warn <utilisateur> [raison]`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?warn @user`\n" +
		"- `?warn @user langage inapproprié`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
