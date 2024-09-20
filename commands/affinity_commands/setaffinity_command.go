package affinity_commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// SetAffinityCommand permet de définir une quantité spécifique d'affinité pour l'auteur
func SetAffinityCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier les droits d'administrateur
	isAdmin, _ := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
	if !isAdmin {
		return
	}

	// Commande de base : ?setaffinity <quantité>
	prefix := fmt.Sprintf("%ssetaffinity", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?setaffinity <quantité>")
		return
	}

	// Convertir la quantité en entier
	affinityAmount, err := strconv.Atoi(args[1])
	if err != nil || affinityAmount < 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide d'affinité.")
		return
	}

	// Définir l'affinité pour l'utilisateur
	err = services.NewUserService(controllers.NewUserController()).SetAffinity(m.Author.ID, affinityAmount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la définition de l'affinité.")
		return
	}

	// Confirmation à l'utilisateur
	response := fmt.Sprintf("Votre affinité a été définie à %d.", affinityAmount)
	s.ChannelMessageSend(m.ChannelID, response)
}
