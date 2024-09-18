package riot_games

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Luthor91/Tenshi/config" // Assurez-vous de remplacer par le chemin correct
)

// Helper function for making API requests
func makeRequest(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Riot-Token", config.AppConfig.RiotAPIKey)
	resp, err := config.AppConfig.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// Champion Rotations
func GetChampionRotations() (interface{}, error) {
	url := fmt.Sprintf("%s/lol/platform/v3/champion-rotations", config.AppConfig.RiotBaseURL)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Summoner by Account ID
func GetSummonerByAccountId(encryptedAccountId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-account/%s", config.AppConfig.RiotBaseURL, encryptedAccountId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Summoner by Summoner ID
func GetSummonerById(encryptedSummonerId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", config.AppConfig.RiotBaseURL, encryptedSummonerId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Summoner by PUUID
func GetSummonerByPuuid(encryptedPUUID string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/summoner/v4/summoners/by-puuid/%s", config.AppConfig.RiotBaseURL, encryptedPUUID)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenger Leagues by Queue
func GetChallengerLeagues(queue string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league/v4/challengerleagues/by-queue/%s", config.AppConfig.RiotBaseURL, queue)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Master Leagues by Queue
func GetMasterLeagues(queue string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league/v4/masterleagues/by-queue/%s", config.AppConfig.RiotBaseURL, queue)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Grandmaster Leagues by Queue
func GetGrandmasterLeagues(queue string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league/v4/grandmasterleagues/by-queue/%s", config.AppConfig.RiotBaseURL, queue)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// League Entries by Summoner ID
func GetLeagueEntriesBySummonerId(encryptedSummonerId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league/v4/entries/by-summoner/%s", config.AppConfig.RiotBaseURL, encryptedSummonerId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// League Entries by Queue, Tier, and Division
func GetLeagueEntries(queue, tier, division string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league/v4/entries/%s/%s/%s", config.AppConfig.RiotBaseURL, queue, tier, division)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// League Exp Entries by Queue, Tier, and Division
func GetLeagueExpEntries(queue, tier, division string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/league-exp/v4/entries/%s/%s/%s", config.AppConfig.RiotBaseURL, queue, tier, division)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Clash Players by Summoner ID
func GetClashPlayersBySummonerId(summonerId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/clash/v1/players/by-summoner/%s", config.AppConfig.RiotBaseURL, summonerId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Clash Teams by Team ID
func GetClashTeamsByTeamId(teamId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/clash/v1/teams/%s", config.AppConfig.RiotBaseURL, teamId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Clash Tournaments by Tournament ID
func GetClashTournamentsById(tournamentId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/clash/v1/tournaments/%s", config.AppConfig.RiotBaseURL, tournamentId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Account by Riot ID
func GetAccountByRiotId(gameName, tagLine string) (interface{}, error) {
	url := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s", config.AppConfig.RiotBaseURL, gameName, tagLine)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Account by PUUID
func GetAccountByPuuid(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/riot/account/v1/accounts/by-puuid/%s", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Platform Data
func GetPlatformData() (interface{}, error) {
	url := fmt.Sprintf("%s/lol/status/v4/platform-data", config.AppConfig.RiotBaseURL)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Match by Match ID
func GetMatchById(matchId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/%s", config.AppConfig.RiotBaseURL, matchId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Matches by PUUID
func GetMatchesByPuuid(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Match Timeline by Match ID
func GetMatchTimelineById(matchId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/match/v5/matches/%s/timeline", config.AppConfig.RiotBaseURL, matchId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenges Percentiles
func GetChallengesPercentiles() (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/challenges/percentiles", config.AppConfig.RiotBaseURL)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenge Leaderboards by Level
func GetChallengeLeaderboards(challengeId string, level string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/challenges/%s/leaderboards/by-level/%s", config.AppConfig.RiotBaseURL, challengeId, level)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenge Percentiles by Challenge ID
func GetChallengePercentiles(challengeId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/challenges/%s/percentiles", config.AppConfig.RiotBaseURL, challengeId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenge Config by Challenge ID
func GetChallengeConfig(challengeId string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/challenges/%s/config", config.AppConfig.RiotBaseURL, challengeId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Player Data by PUUID
func GetPlayerData(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/player-data/%s", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Challenges Config
func GetChallengesConfig() (interface{}, error) {
	url := fmt.Sprintf("%s/lol/challenges/v1/challenges/config", config.AppConfig.RiotBaseURL)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Champion Masteries by PUUID
func GetChampionMasteriesByPuuid(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/champion-mastery/v4/champion-masteries/by-puuid/%s", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Champion Masteries by PUUID and Champion ID
func GetChampionMasteriesByPuuidAndChampionId(puuid string, championId int) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/champion-mastery/v4/champion-masteries/by-puuid/%s/by-champion/%d", config.AppConfig.RiotBaseURL, puuid, championId)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Champion Mastery Scores by PUUID
func GetChampionMasteryScoresByPuuid(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/champion-mastery/v4/scores/by-puuid/%s", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Top Champion Masteries by PUUID
func GetTopChampionMasteriesByPuuid(puuid string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/champion-mastery/v4/champion-masteries/by-puuid/%s/top", config.AppConfig.RiotBaseURL, puuid)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Featured Games
func GetFeaturedGames() (interface{}, error) {
	url := fmt.Sprintf("%s/lol/spectator/v5/featured-games", config.AppConfig.RiotBaseURL)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// Active Games by Summoner PUUID
func GetActiveGamesBySummonerId(encryptedPUUID string) (interface{}, error) {
	url := fmt.Sprintf("%s/lol/spectator/v5/active-games/by-summoner/%s", config.AppConfig.RiotBaseURL, encryptedPUUID)
	var data interface{}
	err := makeRequest(url, &data)
	return data, err
}

// GetChampionData récupère les données des champions depuis l'API de Riot Games
func GetChampionData() (map[string]ChampionDataExtended, error) {
	url := fmt.Sprintf("%s/lol/static-data/v3/champions?api_key=%s", config.AppConfig.RiotBaseURL, config.AppConfig.RiotAPIKey)

	var data ChampionDataResponse
	err := makeRequest(url, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
