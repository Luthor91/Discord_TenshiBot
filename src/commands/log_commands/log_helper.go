package log_commands

import "github.com/bwmarrin/discordgo"

// Afficher les arguments possibles pour la commande ?logs
// Afficher les arguments possibles pour la commande ?logs
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : consultation des logs (`?logs`)**\n\n" +

		"**Usage général :**\n" +
		"`?logs [options] <nombre>`\n\n" +

		"**Paramètre obligatoire :**\n" +
		"- `<nombre>` : nombre de logs à récupérer\n\n" +

		"**Options disponibles :**\n" +
		"- `-n <utilisateur>` : filtrer les logs par utilisateur (mention ou ID)\n" +
		"- `-c <salon>`       : filtrer les logs par salon (mention ou ID)\n" +
		"- `-h`               : afficher cette aide\n\n" +

		"**Comportement :**\n" +
		"- Sans option : affiche les derniers logs globaux\n" +
		"- Avec `-n` : affiche les logs de l'utilisateur\n" +
		"- Avec `-c` : affiche les logs du salon\n" +
		"- Avec `-n` et `-c` : affiche les logs de l'utilisateur dans le salon\n\n" +

		"**Exemples :**\n" +
		"- `?logs 10`\n" +
		"- `?logs -n @Utilisateur 20`\n" +
		"- `?logs -c #general 15`\n" +
		"- `?logs -n @Utilisateur -c #general 5`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}
