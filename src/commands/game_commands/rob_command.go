package game_commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// RobCommand permet de voler de l'argent à un utilisateur
func RobCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%srob", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments de la commande
	args := strings.Fields(m.Content[len(command):])
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur : Vous devez spécifier un utilisateur à voler. Syntaxe : `%s <pseudo_utilisateur>` ou `%s <@mention>`", command, command))
		return
	}

	var targetUser *discordgo.User
	var err error

	// Vérifier si un utilisateur a été mentionné
	if len(m.Mentions) > 0 {
		// Récupérer l'utilisateur mentionné
		targetUser = m.Mentions[0]
	} else {
		// Sinon, rechercher par pseudo
		targetUsername := args[0]
		targetUser, err = discord.FindUserByUsername(s, m.GuildID, targetUsername)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur : Utilisateur `%s` non trouvé.", targetUsername))
			return
		}
	}

	// Récupérer les informations de l'utilisateur cible
	target, err := services.NewUserService().GetUserByDiscordID(targetUser.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur : Impossible de récupérer les informations de l'utilisateur `%s`.", targetUser.Username))
		return
	}

	robAmount := int(float64(target.Money) * 0.05)

	if rand.Intn(100) < 50 { // Taux de réussite de 50%
		user, _ := services.NewUserService().GetUserByDiscordID(m.Author.ID)
		services.NewUserService().AddMoney(target.UserDiscordID, -robAmount)
		services.NewUserService().AddMoney(user.UserDiscordID, robAmount)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Succès ! Vous avez volé %d pièces à %s.", robAmount, targetUser.Username))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Échec ! Vous avez échoué à voler de l'argent à %s.", targetUser.Username))
	}
}
