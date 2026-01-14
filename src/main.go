package main

import (
	"github.com/Luthor91/DiscordBot/commands/banword_commands"
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/core"
	"github.com/Luthor91/DiscordBot/database"
	"github.com/Luthor91/DiscordBot/services"
)

func main() {
	config.LoadConfig(false)
	database.InitDatabase()
	config.CheckConfig()

	wordService := services.NewWordService(database.DB)
	banword_commands.InitWordCommands(wordService)
	
	core.Run()

}
