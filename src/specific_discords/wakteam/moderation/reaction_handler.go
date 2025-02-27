package moderation

import (
	"fmt"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/config"
	"github.com/bwmarrin/discordgo"
)

// OnReactionOnTicket gère les réactions des administrateurs et modérateurs
func OnReactionOnTicket(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Vérifie que la réaction est dans le bon serveur et salon
	if r.GuildID != config.WakteamGuildID || r.ChannelID != config.TicketChannelID {
		return
	}

	// On s'assure que le bot ne réagisse pas a ses propres reactions
	if r.UserID == s.State.User.ID {
		return
	}

	// Vérifie que l'utilisateur est un administrateur ou un modérateur
	isAdmin, errAdmin := discord.UserHasAdminRole(s, r.GuildID, r.UserID)
	isMod, errMod := discord.UserHasModeratorRole(s, r.GuildID, r.UserID)

	// Si l'utilisateur n'est ni admin ni modérateur, on ignore la réaction
	if (errAdmin != nil || !isAdmin) && (errMod != nil || !isMod) {
		fmt.Printf("[INFO] Réaction ignorée : %s n'est ni administrateur ni modérateur.\n", r.UserID)
		return
	}

	// Récupère l'auteur du message d'origine
	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Println("[ERREUR] Impossible de récupérer le message :", err)
		return
	}

	// Vérifie si l'auteur est un utilisateur normal (et non un bot/modérateur)
	if msg.Author.Bot {
		return
	}

	// Détermine l'action à effectuer selon la réaction
	switch r.Emoji.Name {
	case "✅":
		fmt.Printf("[INFO] Ticket de %s validé ✅ par %s.\n", msg.Author.Username, r.UserID)
	case "❌":
		fmt.Printf("[INFO] Ticket de %s rejeté ❌ par %s. Suppression du message.\n", msg.Author.Username, r.UserID)
		err := s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		if err != nil {
			fmt.Println("[ERREUR] Impossible de supprimer le message :", err)
		}
	default:
		fmt.Println("[INFO] Réaction non gérée :", r.Emoji.Name)
	}
}
