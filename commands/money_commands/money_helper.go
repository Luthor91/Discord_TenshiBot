package money_commands

import (
	"fmt"
	"strconv"

	"math/rand"
	"time"

	"github.com/Luthor91/Tenshi/api/discord"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Fonction pour afficher un message d'aide
func displayHelpMessage(s *discordgo.Session, channelID string) {
	message := "Utilisation de la commande money:\n" +
		"-n : Mentionner un utilisateur\n" +
		"-r [montant] : Retirer de l'argent\n" +
		"-d : Récompense quotidienne\n" +
		"-m : Afficher votre solde\n" +
		"-g [montant] : Donner de l'argent à un utilisateur mentionné\n" +
		"-s [montant] : Définir l'argent d'un utilisateur (admin uniquement)"
	s.ChannelMessageSend(channelID, message)
}

// Gérer la mention d'un utilisateur
func handleMention(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.User {
	if len(m.Mentions) > 0 {
		return m.Mentions[0]
	} else {
		s.ChannelMessageSend(m.ChannelID, "Veuillez mentionner un utilisateur avec -n.")
		return nil
	}
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

	userMoney, err := services.NewUserService(controllers.NewUserController()).GetMoney(m.Author.ID)
	if err != nil || userMoney < amount {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas assez d'argent pour retirer ce montant.")
		return
	}
	newMoney := userMoney - amount
	services.NewUserService(controllers.NewUserController()).SetMoney(m.Author.ID, newMoney)
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
	canReceive, timeLeft := services.NewUserService(userController).CanReceiveDailyReward(user)
	if !canReceive {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous devez attendre encore %v avant de réclamer la prochaine récompense quotidienne.", timeLeft.Round(time.Minute)))
		return
	}

	rand.Seed(time.Now().UnixNano())   // Initialisation du générateur de nombres aléatoires
	randomAmount := rand.Intn(91) + 10 // Montant aléatoire entre 10 et 100
	services.NewUserService(userController).UpdateDailyMoney(user, randomAmount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez reçu %d unités aujourd'hui !", randomAmount))
}

// Afficher le solde d'un utilisateur
func handleShowMoney(s *discordgo.Session, m *discordgo.MessageCreate) {
	money, err := services.NewUserService(controllers.NewUserController()).GetMoney(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de votre solde.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez %d unités de monnaie.", money))
}

// Donner de l'argent à un utilisateur
func handleGiveMoney(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	targetUser := handleMention(s, m)
	if targetUser == nil {
		return
	}

	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour donner.")
		return
	}

	amount, err := strconv.Atoi(args[2])
	if err != nil || amount <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour donner.")
		return
	}

	services.NewUserService(controllers.NewUserController()).GiveMoney(m.Author.ID, targetUser.ID, amount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez donné %d unités à %s.", amount, targetUser.Username))
}

// Définir l'argent d'un utilisateur (admin seulement)
func handleSetMoney(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	isAdmin, err := discord.UserHasAdminRole(s, m.GuildID, m.Author.ID)
	if err != nil || !isAdmin {
		s.ChannelMessageSend(m.ChannelID, "Vous n'avez pas la permission de définir la monnaie.")
		return
	}

	targetUser := handleMention(s, m)
	if targetUser == nil {
		return
	}

	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.")
		return
	}

	amount, err := strconv.Atoi(args[2])
	if err != nil || amount < 0 {
		s.ChannelMessageSend(m.ChannelID, "Veuillez entrer un montant valide pour définir l'argent.")
		return
	}

	services.NewUserService(controllers.NewUserController()).SetMoney(targetUser.ID, amount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("L'argent de %s a été défini à %d.", targetUser.Username, amount))
}
