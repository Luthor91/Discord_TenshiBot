package config

import (
	"os"
)

var (
	BotToken  string
	BotPrefix string
)

// LoadConfig charge les variables d'environnement nécessaires
func LoadConfig() {
	BotToken = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("PREFIX")
}
