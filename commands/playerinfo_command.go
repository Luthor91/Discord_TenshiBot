package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

const riotAPIBaseURL = "https://euw1.api.riotgames.com/lol/summoner/v4/summoners/by-name/"

// SummonerInfo représente la structure des informations du joueur
type SummonerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

// PlayerInfoCommand affiche les informations sur un joueur LoL
func PlayerInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	commandPrefix := fmt.Sprintf("%splayerinfo", config.AppConfig.BotPrefix)
	if m.Content == commandPrefix {
		// Extraire le nom du joueur à partir du message
		playerName := m.Content[len(commandPrefix)+1:]

		// Construire l'URL de l'API
		apiURL := fmt.Sprintf("%s%s?api_key=%s", riotAPIBaseURL, playerName, config.AppConfig.RiotAPIKey)

		resp, err := http.Get(apiURL)
		if err != nil {
			log.Println("Erreur lors de l'appel à l'API Riot:", err)
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Erreur lors de la lecture de la réponse:", err)
			return
		}

		var summonerInfo SummonerInfo
		err = json.Unmarshal(body, &summonerInfo)
		if err != nil {
			log.Println("Erreur lors du parsing de la réponse:", err)
			return
		}

		response := fmt.Sprintf("Nom: %s\nNiveau: %d\nID Invocateur: %s\n", summonerInfo.Name, summonerInfo.SummonerLevel, summonerInfo.ID)
		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
	}
}
