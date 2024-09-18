package models

import (
	"net/http"

	"github.com/KnutZuidema/golio"
)

type Config struct {
	BotToken        string
	BotPrefix       string
	RiotAPIKey      string
	LoLEsportAPIKey string
	LoLPatchVersion string
	RiotBaseURL     string

	LoLRegion   string
	LoLServer   string
	Client      *http.Client
	GolioClient *golio.Client
}
