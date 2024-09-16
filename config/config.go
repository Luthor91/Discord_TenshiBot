package config

import (
	"os"

	"github.com/Luthor91/Tenshi/models"
)

var AppConfig models.Config

// LoadConfig charge la configuration depuis un fichier JSON
func LoadConfig() {
	AppConfig = models.Config{
		BotToken:  os.Getenv("TOKEN"),
		BotPrefix: os.Getenv("PREFIX"),
	}
}
