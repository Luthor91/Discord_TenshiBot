package main

import (
	"log"
	"os"

	bot "github.com/Luthor91/Tenshi/Bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
	}
	bot.BotToken = os.Getenv("TOKEN")
	bot.Run()
}
