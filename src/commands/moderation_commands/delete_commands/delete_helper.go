package delete_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Afficher les arguments possibles pour la commande ?delete
func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : suppression de messages (`?delete`)**\n\n" +

		"**Usage général :**\n" +
		"`?delete [options]`\n\n" +

		"**Options disponibles :**\n" +
		"- `-d <nombre>` : nombre de messages à supprimer (obligatoire)\n" +
		"- `-n <pseudo>` : supprimer uniquement les messages de cet utilisateur\n" +
		"- `-c <salon>`  : salon cible (nom, ID ou mention). Par défaut : salon courant\n" +
		"- `-v`          : afficher un message de confirmation\n" +
		"- `-h`          : afficher cette aide\n\n" +

		"**Exemples :**\n" +
		"- `?delete -d 10`\n" +
		"- `?delete -d 20 -n user`\n" +
		"- `?delete -d 5 -c #general -v`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}


// resolveUserID résout un pseudo Discord en ID
func resolveUserID(s *discordgo.Session, guildID, username string) (string, error) {
	// Recherche de l'utilisateur dans le serveur
	members, err := s.GuildMembers(guildID, "", 1000)
	if err != nil {
		return "", fmt.Errorf("impossible de récupérer les membres du serveur : %v", err)
	}

	// Recherche de l'utilisateur dans la liste des membres
	for _, member := range members {
		if strings.EqualFold(member.User.Username, username) {
			return member.User.ID, nil
		}
	}

	// Si l'utilisateur n'est pas trouvé, on retourne une erreur
	return "", fmt.Errorf("utilisateur non trouvé : %s", username)
}

// resolveChannelID résout un nom de salon ou un ID en ID de salon
func resolveChannelID(s *discordgo.Session, guildID, channelArg string) (string, error) {
	// Si l'argument commence par <# et se termine par >, c'est un ID de salon sous forme de mention
	if strings.HasPrefix(channelArg, "<#") && strings.HasSuffix(channelArg, ">") {
		return channelArg[2 : len(channelArg)-1], nil
	}

	// Vérifie si l'argument est directement un ID numérique
	if _, err := strconv.ParseUint(channelArg, 10, 64); err == nil {
		// Vérifie que l'ID existe bien sur le serveur
		channel, err := s.Channel(channelArg)
		if err == nil && channel.GuildID == guildID {
			return channelArg, nil
		}
	}

	// Si l'argument est un nom de salon, rechercher l'ID du salon
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return "", fmt.Errorf("impossible de récupérer les salons du serveur : %v", err)
	}

	// Extrait seulement les caractères alphanumériques du nom recherché
	// pour une correspondance basée uniquement sur le texte visible
	var alphaNumOnly = func(input string) string {
		var result strings.Builder
		for _, r := range input {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
				result.WriteRune(r)
			}
		}
		return result.String()
	}

	cleanedArg := alphaNumOnly(strings.ToLower(channelArg))
	fmt.Printf("Argument nettoyé: [%s] (len=%d)\n", cleanedArg, len(cleanedArg))

	// Deux passes: correspondance exacte, puis correspondance partielle
	for _, channel := range channels {
		// Vérification par texte visible uniquement (en ignorant les caractères spéciaux)
		cleanChannelName := alphaNumOnly(strings.ToLower(channel.Name))

		// Correspondance exacte sur la partie alphanumérique
		if cleanChannelName == cleanedArg {
			fmt.Printf("Correspondance exacte trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	// Deuxième passe: correspondance partielle
	for _, channel := range channels {
		cleanChannelName := alphaNumOnly(strings.ToLower(channel.Name))

		// Correspondance partielle
		if strings.Contains(cleanChannelName, cleanedArg) ||
			strings.Contains(cleanedArg, cleanChannelName) {
			fmt.Printf("Correspondance partielle trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	// Dernière tentative: correspondance phonétique ou visuelle proche
	for _, channel := range channels {
		// Affiche les noms avec leurs codes Unicode pour le débogage
		fmt.Printf("Canal: [%s] -> alphaNum: [%s]\n",
			channel.Name,
			alphaNumOnly(channel.Name))

		// Si le canal contient "ticket" quelque part et que l'argument contient aussi "ticket"
		if strings.Contains(strings.ToLower(channel.Name), "ticket") &&
			strings.Contains(strings.ToLower(channelArg), "ticket") {
			fmt.Printf("Correspondance par mot-clé trouvée: %s\n", channel.Name)
			return channel.ID, nil
		}
	}

	return "", fmt.Errorf("salon non trouvé : %s", channelArg)
}