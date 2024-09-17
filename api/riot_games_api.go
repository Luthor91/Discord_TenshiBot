package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Luthor91/Tenshi/config"
)

// Client structure to hold the base URL and the API key
type RiotAPIClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewRiotAPIClient creates a new instance of RiotAPIClient
func NewRiotAPIClient() *RiotAPIClient {
	return &RiotAPIClient{
		BaseURL: "https://euw1.api.riotgames.com",
		APIKey:  config.AppConfig.RiotAPIKey,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Helper function to perform GET requests
func (client *RiotAPIClient) get(endpoint string, result interface{}) error {
	url := fmt.Sprintf("%s%s", client.BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Riot-Token", client.APIKey)
	resp, err := client.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("API error: %s", body)
	}

	return json.Unmarshal(body, &result)
}

// Account Endpoints
func (client *RiotAPIClient) GetAccountByPuuid(puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/riot/account/v1/accounts/by-puuid/%s", puuid), &result)
	return result, err
}

func (client *RiotAPIClient) GetAccountByRiotID(gameName, tagLine string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/riot/account/v1/accounts/by-riot-id/%s/%s", gameName, tagLine), &result)
	return result, err
}

func (client *RiotAPIClient) GetAccountByAccessToken() (interface{}, error) {
	var result interface{}
	err := client.get("/riot/account/v1/accounts/me", &result)
	return result, err
}

func (client *RiotAPIClient) GetActiveShard(game, puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/riot/account/v1/active-shards/by-game/%s/by-puuid/%s", game, puuid), &result)
	return result, err
}

// Champion Mastery Endpoints
func (client *RiotAPIClient) GetAllChampionMasteries(puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/champion-mastery/v4/champion-masteries/by-puuid/%s", puuid), &result)
	return result, err
}

func (client *RiotAPIClient) GetChampionMasteryByChampion(puuid string, championID int) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/champion-mastery/v4/champion-masteries/by-puuid/%s/by-champion/%d", puuid, championID), &result)
	return result, err
}

func (client *RiotAPIClient) GetTopChampionMasteries(puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/champion-mastery/v4/champion-masteries/by-puuid/%s/top", puuid), &result)
	return result, err
}

func (client *RiotAPIClient) GetChampionMasteryScore(puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/champion-mastery/v4/scores/by-puuid/%s", puuid), &result)
	return result, err
}

// Champion Rotations
func (client *RiotAPIClient) GetChampionRotations() (interface{}, error) {
	var result interface{}
	err := client.get("/lol/platform/v3/champion-rotations", &result)
	return result, err
}

// League Endpoints
func (client *RiotAPIClient) GetChallengerLeague(queue string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/challengerleagues/by-queue/%s", queue), &result)
	return result, err
}

func (client *RiotAPIClient) GetLeagueEntriesBySummoner(summonerID string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/entries/by-summoner/%s", summonerID), &result)
	return result, err
}

func (client *RiotAPIClient) GetLeagueEntries(queue, tier, division string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s", queue, tier, division), &result)
	return result, err
}

func (client *RiotAPIClient) GetGrandmasterLeague(queue string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/grandmasterleagues/by-queue/%s", queue), &result)
	return result, err
}

func (client *RiotAPIClient) GetLeagueByID(leagueID string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/leagues/%s", leagueID), &result)
	return result, err
}

func (client *RiotAPIClient) GetMasterLeague(queue string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/league/v4/masterleagues/by-queue/%s", queue), &result)
	return result, err
}

// Lol Challenges Endpoints
func (client *RiotAPIClient) GetChallengesConfig() (interface{}, error) {
	var result interface{}
	err := client.get("/lol/challenges/v1/challenges/config", &result)
	return result, err
}

func (client *RiotAPIClient) GetChallengePercentiles() (interface{}, error) {
	var result interface{}
	err := client.get("/lol/challenges/v1/challenges/percentiles", &result)
	return result, err
}

func (client *RiotAPIClient) GetChallengeConfig(challengeID int) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/challenges/v1/challenges/%d/config", challengeID), &result)
	return result, err
}

func (client *RiotAPIClient) GetChallengeLeaderboards(challengeID int, level string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/challenges/v1/challenges/%d/leaderboards/by-level/%s", challengeID, level), &result)
	return result, err
}

func (client *RiotAPIClient) GetPlayerChallengeData(puuid string) (interface{}, error) {
	var result interface{}
	err := client.get(fmt.Sprintf("/lol/challenges/v1/player-data/%s", puuid), &result)
	return result, err
}

// LoL Status
func (client *RiotAPIClient) GetPlatformStatus() (interface{}, error) {
	var result interface{}
	err := client.get("/lol/status/v4/platform-data", &result)
	return result, err
}
