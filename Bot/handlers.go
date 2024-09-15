package bot

import (
	"github.com/Luthor91/Tenshi/commands"

	"github.com/bwmarrin/discordgo"
)

// RegisterHandlers enregistre les différentes commandes du bot
func RegisterHandlers(discord *discordgo.Session) {
	discord.AddHandler(commands.BanCommand)
	discord.AddHandler(commands.ByeCommand)
	discord.AddHandler(commands.CalculateCommand)
	discord.AddHandler(commands.DeleteCommand)
	discord.AddHandler(commands.HelpCommand)
	discord.AddHandler(commands.KickCommand)
	discord.AddHandler(commands.PingCommand)
	discord.AddHandler(commands.RandomCommand)

}
