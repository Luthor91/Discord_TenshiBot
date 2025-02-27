package scans

import (
	"log"
	"strings"

	"github.com/Luthor91/DiscordBot/api/discord"
	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/config"

	"github.com/bwmarrin/discordgo"
)

// ScanAndProcessTicketChannel gère la suppression et l'ajout de réactions sur les messages du salon des tickets.
func ScanAndProcessTicketChannel(s *discordgo.Session) {
	channelID := config.TicketChannelID
	guildID := config.WakteamGuildID

	log.Println("[INFO] Début du scan du salon des tickets...")

	// Récupérer tous les messages du salon
	messages, err := s.ChannelMessages(channelID, 50, "", "", "")
	if err != nil {
		log.Println("[ERREUR] Impossible de récupérer les messages du salon :", err)
		return
	}

	log.Printf("[INFO] %d messages récupérés.\n", len(messages))

	for _, msg := range messages {
		log.Printf("[INFO] Traitement du message ID: %s | Auteur: %s | Contenu: %s\n", msg.ID, msg.Author.ID, msg.Content)

		// Vérification des rôles de l'auteur du message
		isAdmin, errAdmin := discord.UserHasAdminRole(s, guildID, msg.Author.ID)
		if errAdmin != nil {
			log.Printf("[ERREUR] Impossible de vérifier si l'utilisateur %s est admin: %v\n", msg.Author.ID, errAdmin)
			continue
		}

		isMod, errMod := discord.UserHasModeratorRole(s, guildID, msg.Author.ID)
		if errMod != nil {
			log.Printf("[ERREUR] Impossible de vérifier si l'utilisateur %s est modérateur: %v\n", msg.Author.ID, errMod)
			continue
		}

		// Si l'auteur est un modérateur ou admin, on ne supprime pas son message
		if isAdmin || isMod {
			log.Printf("[INFO] Message ignoré (Mod/Admin) | Auteur: %s | Contenu: %s\n", msg.Author.ID, msg.Content)
			//continue
		}

		// Vérifier si le message commence par "?report"
		if !strings.HasPrefix(msg.Content, "?report") {
			log.Printf("[INFO] Suppression du message ID: %s | Auteur: %s | Contenu: %s\n", msg.ID, msg.Author.ID, msg.Content)
			err := s.ChannelMessageDelete(channelID, msg.ID)
			if err != nil {
				log.Printf("[ERREUR] Échec de la suppression du message ID: %s | Erreur: %v\n", msg.ID, err)
			} else {
				log.Printf("[SUCCÈS] Message supprimé ID: %s | Auteur: %s\n", msg.ID, msg.Author.ID)
			}
			continue
		}

		// Ajout des réactions ✅ et ❌ si elles ne sont pas déjà présentes
		handleReactions(s, channelID, msg.ID)
	}

	log.Println("[INFO] Fin du scan du salon des tickets.")
}

// handleReactions ajoute ✅ et ❌ aux messages si elles ne sont pas déjà présentes.
func handleReactions(s *discordgo.Session, channelID, messageID string) {
	reactions, err := s.MessageReactions(channelID, messageID, "✅", 100, "", "")
	if err != nil {
		log.Printf("[ERREUR] Impossible de récupérer les réactions ✅ pour le message ID: %s | Erreur: %v\n", messageID, err)
		return
	}
	hasCheck := containsBotReaction(reactions, s.State.User.ID)

	reactions, err = s.MessageReactions(channelID, messageID, "❌", 100, "", "")
	if err != nil {
		log.Printf("[ERREUR] Impossible de récupérer les réactions ❌ pour le message ID: %s | Erreur: %v\n", messageID, err)
		return
	}
	hasCross := containsBotReaction(reactions, s.State.User.ID)

	if !hasCheck {
		log.Printf("[INFO] Ajout de la réaction ✅ au message ID: %s\n", messageID)
		err := s.MessageReactionAdd(channelID, messageID, "✅")
		if err != nil {
			log.Printf("[ERREUR] Échec de l'ajout de la réaction ✅ au message ID: %s | Erreur: %v\n", messageID, err)
		}
	}

	if !hasCross {
		log.Printf("[INFO] Ajout de la réaction ❌ au message ID: %s\n", messageID)
		err := s.MessageReactionAdd(channelID, messageID, "❌")
		if err != nil {
			log.Printf("[ERREUR] Échec de l'ajout de la réaction ❌ au message ID: %s | Erreur: %v\n", messageID, err)
		}
	}
}

// containsBotReaction vérifie si le bot a déjà ajouté une réaction spécifique
func containsBotReaction(reactions []*discordgo.User, botID string) bool {
	for _, user := range reactions {
		if user.ID == botID {
			return true
		}
	}
	return false
}
