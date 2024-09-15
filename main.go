package main

import (
	"log"

	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"

	"github.com/joho/godotenv"
)

func main() {
	// Charger le fichier .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
	}

	config.LoadConfig()

	// Vérifier que le token est présent
	if config.BotToken == "" {
		log.Fatal("Le token du bot est manquant dans le fichier .env")
	}

	// Vérifier que le token est présent
	if config.BotPrefix == "" {
		log.Fatal("Le prefix du bot est manquant dans le fichier .env")
	}

	// Démarrer le bot
	bot.Run()
}
