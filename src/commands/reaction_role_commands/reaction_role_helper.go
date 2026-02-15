package reaction_role_commands

import (
	"strings"

	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)


func findRoleIDByName(s *discordgo.Session, guildID, roleName string) string {

	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return ""
	}

	for _, role := range roles {
		if strings.EqualFold(role.Name, roleName) {
			return role.ID
		}
	}

	return ""
}

func OnReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	if r.UserID == s.State.User.ID {
		return
	}

	service := services.NewReactionRoleService(s)
	service.HandleReactionAdd(r)
}



func showHelpMessage(s *discordgo.Session, channelID string) {
	helpMessage := "" +
		"**Commande : gestion des rôles par réaction (`?rr`)**\n\n" +

		"**Usage général :**\n" +
		"`?rr \"message\" \"reaction1\" \"role1\" [\"reaction2\" \"role2\" ...]`\n\n" +

		"**Paramètres obligatoires :**\n" +
		"- `\"message\"`     : texte que le bot affichera dans le salon\n" +
		"- `\"reaction\"`    : emoji utilisé pour attribuer le rôle\n" +
		"- `\"role\"`        : nom exact du rôle à attribuer\n\n" +

		"**Fonctionnement :**\n" +
		"- Supprime le message de commande\n" +
		"- Publie le message fourni\n" +
		"- Ajoute automatiquement les réactions indiquées\n" +
		"- Ajoute le rôle lorsqu'un utilisateur réagit\n" +
		"- Retire le rôle si la réaction est supprimée\n\n" +

		"**Contraintes :**\n" +
		"- Les réactions et rôles doivent être fournis par paires\n" +
		"- Le bot doit avoir la permission de gérer les rôles\n" +
		"- Le rôle du bot doit être supérieur aux rôles attribués\n\n" +

		"**Exemples :**\n" +
		"- `?rr \"Choisis ton camp\" \"🔥\" \"Feu\" \"❄️\" \"Glace\"`\n`\n"

	s.ChannelMessageSend(channelID, helpMessage)
}

