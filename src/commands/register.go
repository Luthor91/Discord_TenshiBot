package commands

import (
	word_commands "github.com/Luthor91/DiscordBot/commands/banword_commands"
	"github.com/Luthor91/DiscordBot/commands/channel_commands"
	"github.com/Luthor91/DiscordBot/commands/log_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/ban_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/deafen_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/delete_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/kick_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/move_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/mute_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/timeout_commands"
	"github.com/Luthor91/DiscordBot/commands/moderation_commands/warn_commands"
	"github.com/Luthor91/DiscordBot/commands/reaction_role_commands"
	"github.com/Luthor91/DiscordBot/commands/stat_commands"
	"github.com/Luthor91/DiscordBot/commands/utility_commands"
	"github.com/Luthor91/DiscordBot/services"
	"github.com/Luthor91/DiscordBot/specific_discords/wakteam/moderation"

	"github.com/bwmarrin/discordgo"
)

// RegisterHandlers enregistre les différentes commandes du bot
func RegisterHandlers(discord *discordgo.Session) {
	// Créez les services nécessaires
	userServices := services.NewUserService()
	logService := services.NewLogService()

	// Créez le service de message avec les services requis
	messageService := services.NewMessageService(userServices, logService)

	// Enregistrez les gestionnaires de messages
	discord.AddHandler(messageService.NewPrivateMessage)
	discord.AddHandler(messageService.NewServerMessage)

	// Commandes de modération
	discord.AddHandler(delete_commands.DeleteMessageCommand)
	discord.AddHandler(warn_commands.WarnCommand)
	discord.AddHandler(mute_commands.MuteCommand)
	discord.AddHandler(deafen_commands.DeafenCommand)
	discord.AddHandler(move_commands.MoveCommand)
	discord.AddHandler(timeout_commands.TimeoutCommand)
	discord.AddHandler(kick_commands.KickCommand)
	discord.AddHandler(ban_commands.BanCommand)
	

	// Commandes de gestion des mots
	discord.AddHandler(word_commands.WordCommand)

	// Commandes de gestion de channels
	discord.AddHandler(channel_commands.ChannelCommand)

	// Commandes de classement
	// discord.AddHandler(ranking_commands.LeaderboardCommand)
	// discord.AddHandler(ranking_commands.RankCommand)

	// Commandes utilitaires
	discord.AddHandler(utility_commands.ByeCommand)
	discord.AddHandler(utility_commands.HelpCommand)
	discord.AddHandler(utility_commands.RandomCommand)
	discord.AddHandler(utility_commands.CalculateCommand)
	discord.AddHandler(utility_commands.PingCommand)
	discord.AddHandler(utility_commands.ReminderCommand)

	// Commandes de Logs
	discord.AddHandler(log_commands.LogsCommand)

	// Commandes de Stats
	discord.AddHandler(stat_commands.StatCommand)

	// Handler spécifiques, Wakteam
	discord.AddHandler(reaction_role_commands.ReactionRoleCommand)
	discord.AddHandler(reaction_role_commands.OnReactionAdd)


	discord.AddHandler(moderation.OnMessageInTicketChannel)
	discord.AddHandler(moderation.OnReactionOnTicket)

}
