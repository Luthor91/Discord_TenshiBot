package bot

import (
	"github.com/Luthor91/Tenshi/commands"
	"github.com/Luthor91/Tenshi/commands/lol_commands"
	"github.com/Luthor91/Tenshi/services"

	"github.com/bwmarrin/discordgo"
)

// RegisterHandlers enregistre les différentes commandes du bot
func RegisterHandlers(discord *discordgo.Session) {
	discord.AddHandler(services.NewPrivateMessage)
	discord.AddHandler(services.NewServerMessage)

	discord.AddHandler(commands.AffinityCommand)
	discord.AddHandler(commands.AddGoodWordCommand)
	discord.AddHandler(commands.BanCommand)
	discord.AddHandler(commands.ByeCommand)

	discord.AddHandler(commands.CalculateCommand)
	discord.AddHandler(commands.DailyMoneyCommand)
	discord.AddHandler(commands.DeleteCommand)
	discord.AddHandler(commands.ExperienceCommand)
	discord.AddHandler(commands.AddBadWordCommand)
	discord.AddHandler(commands.HelpCommand)
	discord.AddHandler(commands.KickCommand)
	discord.AddHandler(commands.LeaderboardCommand)
	discord.AddHandler(commands.MoneyCommand)
	discord.AddHandler(commands.PingCommand)

	discord.AddHandler(commands.RandomCommand)
	discord.AddHandler(commands.RankCommand)
	discord.AddHandler(commands.ReminderCommand)
	discord.AddHandler(commands.ShopCommand)
	discord.AddHandler(commands.TimeoutCommand)
	discord.AddHandler(commands.UseItemCommand)

	discord.AddHandler(lol_commands.ChampionRotationCommand)

	//discord.AddHandler(lol_commands.SummonerProfileCommand)
	//discord.AddHandler(lol_commands.ChampionInfoCommand)
}
