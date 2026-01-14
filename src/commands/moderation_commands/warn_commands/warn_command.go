package warn_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/bwmarrin/discordgo"
)

// WarnCommand gère la commande ?warn
// Usage : ?warn <user|@mention> [raison]
func WarnCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérification des droits modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isMod {
		return
	}

	command := config.AppConfig.BotPrefix + "warn"

	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// ?warn
	if m.Content == command {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Découpage simple par espaces
	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		showHelpMessage(s, m.ChannelID)
		return
	}

	// Target : un seul mot
	targetRaw := parts[1]
	targetUser := discord.HandleTarget(s, m, targetRaw)
	if targetUser == nil {
		s.ChannelMessageSend(
			m.ChannelID,
			"Utilisateur introuvable.",
		)
		return
	}

	// Raison optionnelle
	reason := "Aucune raison spécifiée"
	if len(parts) > 2 {
		reason = strings.Join(parts[2:], " ")
	}

	// Envoi du warn
	warnService := services.NewWarnService(s, m.GuildID)
	warnCount, err := warnService.AddWarn(
		targetUser.ID,
		reason,
		m.Author.ID,
	)
	if err != nil {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprint(err),
		)
		return
	}

	s.ChannelMessageSend(
		m.ChannelID,
		fmt.Sprintf(
			"Avertissement envoyé à %s (%d warn%s) : %s",
			targetUser.Username,
			warnCount,
			func() string {
				if warnCount > 1 {
					return "s"
				}
				return ""
			}(),
			reason,
		),
	)
}
