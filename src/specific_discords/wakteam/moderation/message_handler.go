package moderation

import (
	"fmt"
	"strings"

	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/config"
	"github.com/bwmarrin/discordgo"
)

// OnMessageInTicketChannel gère les nouveaux messages dans le salon des tickets
func OnMessageInTicketChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Vérifier si le service de ticket est correctement initialisé
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérification des paramètres du message
	if m.GuildID != config.WakteamGuildID || m.ChannelID != config.TicketChannelID {
		return
	}

	// Supprimer le message si ce n'est pas un message valide (pas de préfixe ?report)
	if !strings.HasPrefix(m.Content, "?report") {
		if strings.HasPrefix(m.Content, "?") {
			return
		}
		fmt.Println("[INFO] Suppression du message de ", m.Author.Username, " : message non valide")
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Println("[ERREUR] Impossible de supprimer le message :", err)
		}
		return
	}

	// Traiter le message ?report ...
	content := strings.TrimPrefix(m.Content, "?report ")
	fmt.Println("[DEBUG] Création du ticket pour :", m.Author.Username, "avec le contenu :", content)

	// Au lieu d'enregistrer dans une base de données, simuler la création d'un ticket
	// Vous pouvez remplacer cette ligne par votre logique propre si nécessaire
	fmt.Printf("[SUCCÈS] Ticket créé pour %s avec le contenu : %s\n", m.Author.Username, content)

	// Ajout des réactions ✅ et ❌ au message
	if err := addReactions(s, m); err != nil {
		fmt.Println("[ERREUR] Impossible d'ajouter des réactions :", err)
	}
}

// Ajoute les réactions ✅ et ❌ au message
func addReactions(s *discordgo.Session, m *discordgo.MessageCreate) error {
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		return fmt.Errorf("impossible d'ajouter la réaction ✅ : %w", err)
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
	if err != nil {
		return fmt.Errorf("impossible d'ajouter la réaction ❌ : %w", err)
	}

	return nil
}
