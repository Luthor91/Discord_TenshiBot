package daily_commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"
	"github.com/bwmarrin/discordgo"
)

// Récompense quotidienne
func handleDailyReward(s *discordgo.Session, m *discordgo.MessageCreate) {
	userController := controllers.NewUserController()
	user, err := userController.GetUserByDiscordID(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erreur lors de la récupération de vos informations.")
		return
	}
	canReceive, timeLeft, err := services.NewUserService().CanReceiveDailyReward(user.UserDiscordID)
	if !canReceive || err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous devez attendre encore %v avant de réclamer la prochaine récompense quotidienne.", timeLeft.Round(time.Minute)))
		return
	}

	randomAmount := rand.Intn(91) + 10 // Montant aléatoire entre 10 et 100
	services.NewUserService().UpdateDailyMoney(user.UserDiscordID, randomAmount)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Vous avez reçu %d unités aujourd'hui !", randomAmount))
}
