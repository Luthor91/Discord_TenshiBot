package money_commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Fonction pour afficher un message d'aide
func displayHelpMessage(s *discordgo.Session, channelID string) {
	message := "Utilisation de la commande money:\n" +
		"-n : Utiliser le nom d'utilisateur\n" +
		"-r [montant] : Retirer de l'argent\n" +
		"-d : Récompense quotidienne\n" +
		"-m : Afficher votre solde\n" +
		"-g [montant] : Donner de l'argent à un utilisateur\n" +
		"-s [montant] : Définir l'argent d'un utilisateur (admin uniquement)"
	s.ChannelMessageSend(channelID, message)
}

// Gérer le ciblage d'un utilisateur par son nom
func handleTarget(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.User {
	args := strings.Fields(m.Content)

	if len(args) < 3 || args[1] != "-n" {
		s.ChannelMessageSend(m.ChannelID, "Veuillez utiliser la commande avec -n suivi du nom d'utilisateur.")
		return nil
	}

	target := args[2]

	// Vérification des utilisateurs dans les mentions
	for _, mention := range m.Mentions {
		if mention.Username == target || fmt.Sprintf("%s#%s", mention.Username, mention.Discriminator) == target {
			return mention
		}
	}

	users, err := s.GuildMembers(m.GuildID, "", 100)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération des membres.")
		return nil
	}

	for _, user := range users {
		if user.User.Username == target || fmt.Sprintf("%s#%s", user.User.Username, user.User.Discriminator) == target {
			return user.User
		}
	}

	s.ChannelMessageSend(m.ChannelID, "Utilisateur non trouvé.")
	return nil
}

// Retirer de l'argent à un utilisateur
func handleRemoveMoney(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour retirer.")
		return
	}
	amount, err := strconv.Atoi(args[2])
	if err != nil || amount <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour retirer.")
		return
	}

	userMoney, err := services.NewUserService().GetMoney(m.Author.ID)
	if err != nil || userMoney < amount {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour retirer ce montant.")
		return
	}
	newMoney := userMoney - amount
	services.NewUserService().SetMoney(m.Author.ID, newMoney)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez retiré %d unités. Nouveau solde : %d", amount, newMoney))
}

// Récompense quotidienne
func handleDailyReward(s *discordgo.Session, m *discordgo.MessageCreate) {
	userController := controllers.NewUserController()
	user, err := userController.GetUserByDiscordID(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de vos informations.")
		return
	}
	canReceive, timeLeft, err := services.NewUserService().CanReceiveDailyReward(user)
	if !canReceive || err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous devez attendre encore %v avant de réclamer la prochaine récompense quotidienne.", timeLeft.Round(time.Minute)))
		return
	}

	randomAmount := rand.Intn(91) + 10
	services.NewUserService().UpdateDailyMoney(user, randomAmount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez reçu %d unités aujourd'hui !", randomAmount))
}

// Afficher le solde d'un utilisateur
func handleShowMoney(s *discordgo.Session, m *discordgo.MessageCreate, targetUser *discordgo.User) {
	money, err := services.NewUserService().GetMoney(targetUser.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération du solde.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s a %d unités de monnaie.", targetUser.Username, money))
}

// Donner de l'argent à un utilisateur
func handleGiveMoney(s *discordgo.Session, m *discordgo.MessageCreate, args []string, targetUser *discordgo.User) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour donner.")
		return
	}

	amount, err := strconv.Atoi(args[2])
	if err != nil || amount <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour donner.")
		return
	}

	services.NewUserService().GiveMoney(m.Author.ID, targetUser.ID, amount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez donné %d unités à %s.", amount, targetUser.Username))
}

// Définir l'argent d'un utilisateur (admin seulement)
func handleSetMoney(s *discordgo.Session, m *discordgo.MessageCreate, args []string, targetUser *discordgo.User) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.")
		return
	}

	amount, err := strconv.Atoi(args[2])
	if err != nil || amount < 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.")
		return
	}

	services.NewUserService().SetMoney(targetUser.ID, amount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'argent de %s a été défini à %d.", targetUser.Username, amount))
}

// handleShowUserBalance affiche le solde de l'utilisateur
func handleShowUserBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	balance, _ := services.NewUserService().GetMoney(userID)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, votre solde est de %d.", m.Author.Username, balance))
}
