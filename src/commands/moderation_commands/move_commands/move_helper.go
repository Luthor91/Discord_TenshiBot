package move_commands

import "github.com/bwmarrin/discordgo"

// Afficher l'aide pour la commande ?move
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : déplacement d'un utilisateur (`?move`)**\n\n" +

		"**Usage :**\n" +
		"`?move <utilisateur> <salon>`\n\n" +

		"**Paramètres :**\n" +
		"- `utilisateur` : pseudo ou mention (obligatoire)\n" +
		"- `salon`       : salon vocal cible (nom, ID ou mention)\n\n" +

		"**Exemples :**\n" +
		"- `?move @user Vocal-1`\n" +
		"- `?move user 1234567890`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
