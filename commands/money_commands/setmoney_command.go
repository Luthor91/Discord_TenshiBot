package money_commands

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

// SetMoneyCommand permet de définir une quantité spécifique de monnaie pour l'auteur
func SetMoneyCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Vérifier les droits d'administrateur
	isAdmin, _ := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
	if !isAdmin {
		return
	}

	// Commande de base : ?setmoney <quantité>
	prefix := fmt.Sprintf("%ssetmoney", config.AppConfig.BotPrefix)
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}
	args := strings.Fields(m.Content)
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Utilisation : ?setmoney <quantité>")
		return
	}

	// Convertir la quantité d'argent en entier
	moneyAmount, err := strconv.Atoi(args[1])
	if err != nil || moneyAmount < 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer une quantité valide de monnaie.")
		return
	}

	// Définir la monnaie pour l'utilisateur
	err = services.NewUserService(controllers.NewUserController()).SetMoney(m.Author.ID, moneyAmount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la définition de la monnaie.")
		return
	}

	// Confirmation à l'utilisateur
	response := fmt.Sprintf("Votre monnaie a été définie à %d unités.", moneyAmount)
	s.ChannelMessageSend(m.ChannelID, response)

}
