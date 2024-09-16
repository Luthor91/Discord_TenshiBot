package main

import (
	"log"

	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"
	"github.com/joho/godotenv"
)

func main() {
	// Charger le fichier .env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
	}

	// Charger la configuration depuis le fichier .env
	config.LoadConfig()

	// Vérifier si le token et le préfixe sont définis dans le fichier .env
	if config.AppConfig.BotToken == "" {
		log.Fatal("Le token du bot est manquant dans le fichier .env.")
	}
	if config.AppConfig.BotPrefix == "" {
		log.Fatal("Le préfixe du bot est manquant dans le fichier .env.")
	}

	// Charger les mots (banwords, goodwords)
	features.LoadWords()

	// Démarrer le bot
	bot.Run()
}
