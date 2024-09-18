package main

import (
	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"
)

func main() {
	// Charger la configuration depuis le fichier .env
	config.LoadConfig()
	config.CheckConfig()
	bot.Run()
}
