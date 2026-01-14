package deafen_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?deafen
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : rendre un utilisateur sourd (`?deafen`)**\n\n" +

		"**Usage :**\n" +
		"`?deafen <utilisateur> <durée> [raison]`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `durée`       : durée du deaf (ex: `10m`, `1h`, `2h30`) (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?deafen @user 10m`\n" +
		"- `?deafen @user 1h bruit excessif`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
