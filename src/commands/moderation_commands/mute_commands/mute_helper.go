package mute_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?mute
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : rendre un utilisateur muet (`?mute`)**\n\n" +

		"**Usage :**\n" +
		"`?mute <utilisateur> <durée> [raison]`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `durée`       : durée du mute (ex: `10m`, `1h`, `2h30`) (obligatoire)\n" +
		"- `raison`      : optionnelle (par défaut : \"Aucune raison spécifiée\")\n\n" +

		"**Exemples :**\n" +
		"- `?mute @user 10m`\n" +
		"- `?mute @user 1h spam vocal`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
