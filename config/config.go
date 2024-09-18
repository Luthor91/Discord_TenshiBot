package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/api"
	"github.com/Luthor91/Tenshi/models"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var AppConfig models.Config

// LoadConfig charge la configuration depuis un fichier JSON
func LoadConfig() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Erreur lors du chargement du fichier .env: %v", err)
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
		Client:          &http.Client{Timeout: 10 * time.Second},
		GolioClient:     client,
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
