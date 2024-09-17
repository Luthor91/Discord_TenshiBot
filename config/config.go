package config

import (
	"log"
	"os"

	"github.com/Luthor91/Tenshi/models"
	"github.com/joho/godotenv"
)

var AppConfig models.Config

// LoadConfig charge la configuration depuis un fichier JSON
func LoadConfig() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
	}

	AppConfig = models.Config{
		BotToken:  os.Getenv("TOKEN"),
		BotPrefix: os.Getenv("PREFIX"),
	}
}

func CheckConfig() {
	// Vérifier si le token et le préfixe sont définis dans le fichier .env
	if AppConfig.BotToken == "" {
		log.Fatal("Le token du bot est manquant dans le fichier .env.")
	}
	if AppConfig.BotPrefix == "" {
		log.Fatal("Le préfixe du bot est manquant dans le fichier .env.")
	}

}
