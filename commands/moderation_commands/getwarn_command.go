package moderation_commands

import (
	"fmt"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// GetWarnCommand récupère les avertissements d'un utilisateur
func GetWarnCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorer les messages du bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier si l'utilisateur est modérateur
	isMod, err := discord.UserHasModeratorRole(s, m.GuildID, m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la vérification des rôles.")
		return
	}
	if !isMod {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas les permissions nécessaires pour récupérer les avertissements.")
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sgetwarn", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Récupérer la mention
	parts := strings.Fields(m.Content)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !getwarn <@mention_user>")
		return
	}

	mention := parts[1]
	if !strings.HasPrefix(mention, "<@") || !strings.HasSuffix(mention, ">") {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur valide.")
		return
	}
	mentionID := mention[2 : len(mention)-1] // Extrait l'ID de la mention (enlève <@ et >)

	// Récupérer les avertissements de l'utilisateur via le service
	warnController := controllers.NewWarnController()
	warnService := services.NewWarnService(warnController, s, m.GuildID)
	warns, err := warnService.GetWarns(mentionID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des avertissements.")
		return
	}

	// Vérifier si des avertissements existent
	if len(warns) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Aucun avertissement trouvé pour <@%s>.", mentionID))
		return
	}

	// Créer un message listant les avertissements
	warnMessages := ""
	for _, warn := range warns {
		warnMessages += fmt.Sprintf("- Raison : %s, Averti par : <@%s>, Date : %s\n", warn.Reason, warn.AdminID, warn.CreatedAt.Format("02 Jan 2006 15:04"))
	}

	// Envoyer la liste des avertissements
	response := fmt.Sprintf("**Avertissements pour <@%s> :**\n%s", mentionID, warnMessages)
	s.ChannelMessageSend(m.ChannelID, response)
}
