package game_commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

func InvestCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Formater la commande avec le préfixe
	command := fmt.Sprintf("%sinvest", config.AppConfig.BotPrefix)

	// Vérifier si le message commence par la commande
	if !strings.HasPrefix(m.Content, command) {
		return
	}

	// Extraire les arguments après la commande
	args := strings.TrimSpace(strings.TrimPrefix(m.Content, command))

	// Vérifier si aucun argument n'est fourni
	if args == "" {
		// Retourner la syntaxe complète de la commande
		syntax := fmt.Sprintf("Syntaxe : %sinvest <montant>", config.AppConfig.BotPrefix)
		s.ChannelMessageSend(m.ChannelID, syntax)
		return
	}

	// Séparer les arguments en fonction des espaces
	argsList := strings.Fields(args)

	// Vérifier qu'il y a assez d'arguments
	if len(argsList) < 1 {
		syntax := fmt.Sprintf("Syntaxe : %sinvest <montant>", config.AppConfig.BotPrefix)
		s.ChannelMessageSend(m.ChannelID, syntax)
		return
	}

	// Convertir le premier argument en un entier (montant)
	amount, err := strconv.Atoi(argsList[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Veuillez fournir un montant valide.")
		return
	}

	// Récupérer les informations de l'utilisateur
	userService := services.NewUserService()
	user, err := userService.GetUserByDiscordID(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la récupération des informations utilisateur : %v", err))
		return
	}

	// Vérifier si l'utilisateur a déjà investi il y a moins de 10 heures
	investmentService := services.NewInvestmentService()
	lastInvestment, err := investmentService.GetLastInvestmentUser(m.Author.ID)
	if err == nil && time.Since(lastInvestment.InvestedAt).Hours() < 10 {
		s.ChannelMessageSend(m.ChannelID, "Vous avez déjà investi au cours des dernières 10 heures. Revenez plus tard.")
		return
	}

	// Vérifier si l'utilisateur a suffisamment d'argent pour investir
	if user.Money < amount {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour effectuer cet investissement.")
		return
	}

	// Déduire l'argent de l'utilisateur
	newMoneyAmount := user.Money - amount
	err = userService.UpdateMoney(m.Author.ID, newMoneyAmount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erreur lors de la mise à jour du solde : %v", err))
		return
	}

	// Enregistrer l'investissement dans la base de données
	err = investmentService.CreateInvestment(m.Author.ID, amount)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de l'enregistrement de l'investissement.")
		return
	}

	// Envoyer un message de confirmation
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez investi %d ! Revenez dans 24 heures pour récupérer votre retour sur investissement.", amount))
}
