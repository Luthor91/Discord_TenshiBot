package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/Luthor91/Tenshi/config"
	"github.com/bwmarrin/discordgo"
)

const riotEsportURL = "https://api.lolesports.com/persisted/gw/getSchedule?hl=en-US"

// EsportMatch représente la structure des matchs
type EsportMatch struct {
	League    string    `json:"league"`
	Team1     string    `json:"team1"`
	Team2     string    `json:"team2"`
	MatchTime time.Time `json:"startTime"`
}

// EsportResponse représente la réponse de l'API
type EsportResponse struct {
	Data struct {
		Schedule struct {
			Events []struct {
				League struct {
					Name string `json:"name"`
				} `json:"league"`
				Match struct {
					Teams []struct {
						Name string `json:"name"`
					} `json:"teams"`
				} `json:"match"`
				StartTime string `json:"startTime"`
			} `json:"events"`
		} `json:"schedule"`
	} `json:"data"`
}

// LoLEsportCommand affiche les matchs de la semaine
func LoLEsportCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	command := fmt.Sprintf("%slolesport", config.AppConfig.BotPrefix)

	if m.Content == command {
		resp, err := http.Get(riotEsportURL)
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

		var esportData EsportResponse
		err = json.Unmarshal(body, &esportData)
		if err != nil {
			log.Println("Erreur lors du parsing de la réponse:", err)
			return
		}

		var response string
		for _, event := range esportData.Data.Schedule.Events {
			startTime, _ := time.Parse(time.RFC3339, event.StartTime)
			response += fmt.Sprintf("Match: %s vs %s\nLeague: %s\nHeure: %s\n\n", event.Match.Teams[0].Name, event.Match.Teams[1].Name, event.League.Name, startTime.Format("Mon, 02 Jan 2006 15:04"))
		}

		_, err = s.ChannelMessageSend(m.ChannelID, response)
		if err != nil {
			log.Println("Erreur lors de l'envoi du message:", err)
		}
	}
}
