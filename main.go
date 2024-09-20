package main

import (
	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/database"
	"github.com/Luthor91/Tenshi/utils"
)

func main() {
	// Charger la configuration depuis le fichier .env
	go utils.PrintMemoryUsage(10)

	config.LoadConfig()
	database.InitDatabase()
	config.CheckConfig()
	bot.Run()

}
