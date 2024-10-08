package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/Luthor91/DiscordBot/models"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var AppConfig models.Config

func LoadConfig(loadEnv bool) {
	if loadEnv {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
		}
	}

	client := golio.NewClient(os.Getenv("API_RIOT"),
		golio.WithRegion(api.RegionEuropeWest),
		golio.WithLogger(logrus.New().WithField("foo", "bar")))

	AppConfig = models.Config{
		BotToken:        os.Getenv("TOKEN"),
		BotPrefix:       os.Getenv("PREFIX"),
		RiotAPIKey:      os.Getenv("API_RIOT"),
		LoLEsportAPIKey: os.Getenv("API_LOL_ESPORT"),
		LoLPatchVersion: os.Getenv("LOL_PATCH_VERSION"),
		RiotBaseURL:     fmt.Sprintf("https://%s.api.riotgames.com", os.Getenv("LOL_REGION")),
		LoLRegion:       os.Getenv("LOL_REGION"),
		LoLServer:       os.Getenv("LOL_SERVER"),

		DBName:          os.Getenv("DB_NAME"),
		DBAdminUser:     os.Getenv("DB_ADMIN_USER"),
		DBAdminPassword: os.Getenv("DB_ADMIN_PASSWORD"),
		DBHost:          os.Getenv("DB_HOST"),
		DBPort:          os.Getenv("DB_PORT"),
		DBSSLMode:       os.Getenv("DB_SSL_MODE"),
		DBUser:          os.Getenv("DB_USER"),
		DBPassword:      os.Getenv("DB_PASSWORD"),

		Client:      &http.Client{Timeout: 10 * time.Second},
		GolioClient: client,
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
