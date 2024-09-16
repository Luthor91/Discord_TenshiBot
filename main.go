package main

import (
	"flag"
	"log"

	"github.com/Luthor91/Tenshi/bot"
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/features"

	"github.com/joho/godotenv"
)

func main() {
	// Définir les variables pour les arguments de ligne de commande
	var token string
	var prefix string

	// Définir les flags pour les arguments de ligne de commande
	flag.StringVar(&token, "token", "", "Le token du bot Discord")
	flag.StringVar(&prefix, "prefix", "", "Le préfixe des commandes du bot")
	flag.Parse()

	// Charger le fichier .env si les arguments ne sont pas fournis
	if token == "" || prefix == "" {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
		}

		// Charger la configuration depuis le fichier .env
		config.LoadConfig()

		// Si les valeurs ne sont toujours pas définies, les signaler comme manquantes
		if config.BotToken == "" {
			log.Fatal("Le token du bot est manquant dans le fichier .env ou en argument.")
		}
		if config.BotPrefix == "" {
			log.Fatal("Le prefix du bot est manquant dans le fichier .env ou en argument.")
		}

		// Utiliser les valeurs du fichier .env
		token = config.BotToken
		prefix = config.BotPrefix
	} else {
		// Utiliser les valeurs des arguments de ligne de commande
		config.BotToken = token
		config.BotPrefix = prefix
	}

	// Charger les mots (banwords, goodwords)
	features.LoadWords()

	// Démarrer le bot
	bot.Run()
}
