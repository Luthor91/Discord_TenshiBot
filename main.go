package main

import (
	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/database"
)

func main() {
	// Charger la configuration depuis le fichier .env
	config.LoadConfig()
	database.InitDatabase()
	config.CheckConfig()
	bot.Run()
}
