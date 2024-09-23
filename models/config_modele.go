package models

import (
	"net/http"

	"github.com/KnutZuidema/golio"
)

// Config contient les configurations de l'application, y compris les paramètres pour la connexion à PostgreSQL
type Config struct {
	// Configuration du bot Discord
	BotToken  string
	BotPrefix string

	// Clés API pour Riot et LoL Esports
	RiotAPIKey      string
	LoLEsportAPIKey string
	LoLPatchVersion string

	// URL de base pour l'API Riot
	RiotBaseURL string

	// Informations pour la connexion à la base de données PostgreSQL
	DBName          string
	DBAdminUser     string
	DBAdminPassword string
	DBHost          string
	DBPort          string
	DBSSLMode       string
	DBUser          string
	DBPassword      string

	// Configuration spécifique à LoL
	LoLRegion string
	LoLServer string

	// Clients
	Client      *http.Client
	GolioClient *golio.Client
}
