package commands

import (
	"github.com/Luthor91/DiscordBot/commands/affinity_commands"
	"github.com/Luthor91/DiscordBot/commands/channel_commands"
	"github.com/Luthor91/DiscordBot/commands/daily_commands"
	"github.com/Luthor91/DiscordBot/commands/experience_commands"
	"github.com/Luthor91/DiscordBot/commands/game_commands"
	"github.com/Luthor91/DiscordBot/commands/item_commands"
	"github.com/Luthor91/DiscordBot/commands/log_commands"
	"github.com/Luthor91/DiscordBot/commands/lol_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands"
	"github.com/Luthor91/DiscordBot/commands/money_commands"
	"github.com/Luthor91/DiscordBot/commands/ranking_commands"
	"github.com/Luthor91/DiscordBot/commands/shop_commands"
	"github.com/Luthor91/DiscordBot/commands/stat_commands"
	"github.com/Luthor91/DiscordBot/commands/utility_commands"
	"github.com/Luthor91/DiscordBot/commands/word_commands"
	"github.com/Luthor91/DiscordBot/services"

	"github.com/bwmarrin/discordgo"
)

// RegisterHandlers enregistre les différentes commandes du bot
func RegisterHandlers(discord *discordgo.Session) {
	// Créez les services nécessaires
	userServices := services.NewUserService()
	affinityService := services.NewAffinityService()
	logService := services.NewLogService()

	// Créez le service de message avec les services requis
	messageService := services.NewMessageService(userServices, affinityService, logService)

	// Enregistrez les gestionnaires de messages
	discord.AddHandler(messageService.NewPrivateMessage)
	discord.AddHandler(messageService.NewServerMessage)

	// Commandes d'affinité
	discord.AddHandler(affinity_commands.AffinityCommand)

	// Commandes de modération
	discord.AddHandler(moderation_commands.ModerateUserCommand)
	discord.AddHandler(moderation_commands.ModerateMessageCommand)

	// Commandes de gestion des mots
	discord.AddHandler(word_commands.WordCommand)

	// Commandes de gestion de channels
	discord.AddHandler(channel_commands.ChannelCommand)

	// Commandes d'expérience
	discord.AddHandler(experience_commands.ExperienceCommand)

	// Commandes de classement
	discord.AddHandler(ranking_commands.LeaderboardCommand)
	discord.AddHandler(ranking_commands.RankCommand)

	// Commandes de monnaie
	discord.AddHandler(money_commands.MoneyCommand)

	discord.AddHandler(shop_commands.ShopCommand)
	discord.AddHandler(shop_commands.BuyCommand)

	// Commandes utilitaires
	discord.AddHandler(utility_commands.ByeCommand)
	discord.AddHandler(utility_commands.HelpCommand)
	discord.AddHandler(utility_commands.RandomCommand)
	discord.AddHandler(utility_commands.CalculateCommand)
	discord.AddHandler(utility_commands.PingCommand)
	discord.AddHandler(utility_commands.ReminderCommand)

	// Commandes ponctuelles
	discord.AddHandler(daily_commands.DailyCommand)

	// Commandes d'articles
	discord.AddHandler(item_commands.ItemCommand)

	// Commandes de Logs
	discord.AddHandler(log_commands.LogsCommand)

	// Commandes de Stats
	discord.AddHandler(stat_commands.StatCommand)

	// Commandes Games
	discord.AddHandler(game_commands.BetCommand)
	discord.AddHandler(game_commands.GuessCommand)
	discord.AddHandler(game_commands.ShifumiCommand)
	discord.AddHandler(game_commands.InvestCommand)
	discord.AddHandler(game_commands.CollectInvestCommand)
	discord.AddHandler(game_commands.RobCommand)

	// Commandes LOL
	discord.AddHandler(lol_commands.ChampionRotationCommand)
	discord.AddHandler(lol_commands.SummonerProfileCommand)
	discord.AddHandler(lol_commands.ChampionInfoCommand)
}
