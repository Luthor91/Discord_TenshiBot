package ban_commands

import "github.com/bwmarrin/discordgo"

// Affiche l'aide pour la commande ?ban
func showHelpMessage(s *discordgo.Session, channelID string) {
	msg := "" +
		"**Commande : bannissement (`?ban`)**\n\n" +
		"Banni définitivement un utilisateur du serveur. Ne peut être débanni que depuis les paramètres du serveur\n\n" +
		"**Usage :**\n" +
		"`?ban <utilisateur> [raison]`\n\n" +

		"**Détails :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?ban @user`\n" +
		"- `?ban @user spam répété`\n"

	s.ChannelMessageSend(channelID, msg)
}
