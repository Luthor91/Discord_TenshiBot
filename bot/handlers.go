package bot

import (
	"github.com/Luthor91/Tenshi/commands/affinity_commands"
	"github.com/Luthor91/Tenshi/commands/channel_commands"
	"github.com/Luthor91/Tenshi/commands/experience_commands"
	"github.com/Luthor91/Tenshi/commands/item_commands"
	"github.com/Luthor91/Tenshi/commands/lol_commands"
	"github.com/Luthor91/Tenshi/commands/moderation_commands"
	"github.com/Luthor91/Tenshi/commands/modvoice_commands"
	"github.com/Luthor91/Tenshi/commands/money_commands"
	"github.com/Luthor91/Tenshi/commands/ranking_commands"
	"github.com/Luthor91/Tenshi/commands/utility_commands"
	"github.com/Luthor91/Tenshi/commands/word_commands"
	"github.com/Luthor91/Tenshi/controllers"
	"github.com/Luthor91/Tenshi/services"

	"github.com/bwmarrin/discordgo"
)

// RegisterHandlers enregistre les différentes commandes du bot
func RegisterHandlers(discord *discordgo.Session) {
	// Créez les services nécessaires
	userServices := services.NewUserService(controllers.NewUserController())
	affinityService := services.NewAffinityService(discord)
	logController := controllers.NewLogController() // Instancier LogController
	logService := services.NewLogService(logController)

	// Créez le service de message avec les services requis
	messageService := services.NewMessageService(userServices, affinityService, logService)

	// Enregistrez les gestionnaires de messages
	discord.AddHandler(messageService.NewPrivateMessage)
	discord.AddHandler(messageService.NewServerMessage)

	// Commandes d'affinité
	discord.AddHandler(affinity_commands.BurnAffinityCommand)
	discord.AddHandler(affinity_commands.GetAffinityCommand)
	discord.AddHandler(affinity_commands.SetAffinityCommand)

	// Commandes de modération
	discord.AddHandler(moderation_commands.BanCommand)
	discord.AddHandler(moderation_commands.KickCommand)
	discord.AddHandler(moderation_commands.DeleteCommand)
	discord.AddHandler(moderation_commands.TimeoutCommand)

	// Commandes de modération vocale
	discord.AddHandler(modvoice_commands.DeafenVoiceCommand)
	discord.AddHandler(modvoice_commands.KickVoiceCommand)
	discord.AddHandler(modvoice_commands.MuteVoiceCommand)
	discord.AddHandler(modvoice_commands.MoveUserVoiceCommand)

	// Commandes de gestion des mots
	discord.AddHandler(word_commands.AddGoodWordCommand)
	discord.AddHandler(word_commands.AddBadWordCommand)
	discord.AddHandler(word_commands.DeleteGoodWordCommand)
	discord.AddHandler(word_commands.DeleteBadWordCommand)
	discord.AddHandler(word_commands.GetGoodWordsCommand)
	discord.AddHandler(word_commands.GetBadWordsCommand)

	// Commandes de gestion de channels
	discord.AddHandler(channel_commands.CreateTextChannelCommand)
	discord.AddHandler(channel_commands.CreateVoiceChannelCommand)
	discord.AddHandler(channel_commands.DeleteChannelByNameCommand)
	discord.AddHandler(channel_commands.CreateTemporaryTextChannelCommand)
	discord.AddHandler(channel_commands.CreateTemporaryVoiceChannelCommand)

	// Commandes d'expérience
	discord.AddHandler(experience_commands.BurnExperienceCommand)
	discord.AddHandler(experience_commands.GetExperienceCommand)
	discord.AddHandler(experience_commands.SetXPCommand)
	discord.AddHandler(experience_commands.GiveXPCommand)

	// Commandes de classement
	discord.AddHandler(ranking_commands.LeaderboardCommand)
	discord.AddHandler(ranking_commands.RankCommand)

	// Commandes de monnaie
	discord.AddHandler(money_commands.BurnMoneyCommand)
	discord.AddHandler(money_commands.GetMoneyCommand)
	discord.AddHandler(money_commands.DailyMoneyCommand)
	discord.AddHandler(money_commands.GiveMoneyCommand)
	discord.AddHandler(money_commands.SetMoneyCommand)

	// Commandes utilitaires
	discord.AddHandler(utility_commands.ByeCommand)
	discord.AddHandler(utility_commands.HelpCommand)
	discord.AddHandler(utility_commands.RandomCommand)
	discord.AddHandler(utility_commands.CalculateCommand)
	discord.AddHandler(utility_commands.PingCommand)
	discord.AddHandler(utility_commands.ReminderCommand)

	// Commandes d'articles
	discord.AddHandler(item_commands.InventoryCommand)
	discord.AddHandler(item_commands.UseItemCommand)
	discord.AddHandler(item_commands.GiveItemCommand)

	// Commandes LOL
	discord.AddHandler(lol_commands.ChampionRotationCommand)
	discord.AddHandler(lol_commands.SummonerProfileCommand)
	discord.AddHandler(lol_commands.ChampionInfoCommand)
}
