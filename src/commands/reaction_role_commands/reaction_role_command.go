package reaction_role_commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

func ReactionRoleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérification modérateur
	isMod, _ := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if !isMod {
		return
	}

	command := fmt.Sprintf("%srr", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Suppression immédiate du message utilisateur
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	// Extraction des arguments entre guillemets
	re := regexp.MustCompile(`"([^"]+)"`)
	matches := re.FindAllStringSubmatch(m.Content, -1)

	if len(matches) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Format invalide. Voir aide.")
		return
	}

	displayMessage := matches[0][1]

	// Vérifier paires emoji / rôle
	if (len(matches)-1)%2 != 0 {
		s.ChannelMessageSend(m.ChannelID, "Les réactions et rôles doivent être en paires.")
		return
	}

	roleMappings := make(map[string]string)

	// Vérification complète des rôles AVANT toute action
	for i := 1; i < len(matches); i += 2 {

		emoji := matches[i][1]
		roleName := matches[i+1][1]

		roleID := findRoleIDByName(s, m.GuildID, roleName)
		if roleID == "" {
			s.ChannelMessageSend(m.ChannelID,
				fmt.Sprintf("Rôle introuvable : %s", roleName),
			)
			return
		}

		roleMappings[emoji] = roleID
	}

	// Envoi du message du bot
	sentMsg, err := s.ChannelMessageSend(m.ChannelID, displayMessage)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'envoi du message.")
		return
	}

	// Ajout des réactions
	for emoji := range roleMappings {
		_ = s.MessageReactionAdd(sentMsg.ChannelID, sentMsg.ID, emoji)
	}

	// Sauvegarde en base via Service
	rrService := services.NewReactionRoleService(s)

	err = rrService.CreateReactionRoles(
		m.GuildID,
		sentMsg.ChannelID,
		sentMsg.ID,
		roleMappings,
	)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la sauvegarde en base.")
		return
	}
}
