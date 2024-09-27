package main

import (
	"github.com/Luthor91/DiscordBot/config"
	"github.com/Luthor91/DiscordBot/core"
	"github.com/Luthor91/DiscordBot/database"
)

func main() {
	config.LoadConfig(false)
	database.InitDatabase()
	config.CheckConfig()
	core.Run()

}
