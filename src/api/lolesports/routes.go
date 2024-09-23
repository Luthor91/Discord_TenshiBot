package lolesports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Luthor91/Tenshi/config"
)

// API Key and URLs
var (
	apiKey       = config.AppConfig.LoLEsportAPIKey
	apiURL       = "https://esports-api.lolesports.com/persisted/gw"
	liveStatsAPI = "https://feed.lolesports.com/livestats/v1"
)

// Response structure for API responses
type APIResponse struct {
	Data json.RawMessage `json:"data"`
}

// LolesportsAPI structure
type LolesportsAPI struct {
	client *http.Client
}

// NewLolesportsAPI creates a new instance of LolesportsAPI
func NewLolesportsAPI() *LolesportsAPI {
	return &LolesportsAPI{
		client: &http.Client{},
	}
}

// getLatestDate returns the current date and time in ISO 8601 format
func getLatestDate() string {
	now := time.Now().UTC().Truncate(time.Second)
	return now.Format(time.RFC3339)
}

// helper function to make GET requests
func (api *LolesportsAPI) getRequest(url string, params map[string]string) (*APIResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", apiKey)
	query := req.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetLeagues retrieves leagues data
func (api *LolesportsAPI) GetLeagues(hl string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getLeagues", map[string]string{"hl": hl})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetTournamentsForLeague retrieves tournaments for a league
func (api *LolesportsAPI) GetTournamentsForLeague(hl string, leagueID string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getTournamentsForLeague", map[string]string{"hl": hl, "leagueId": leagueID})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetStandings retrieves standings
func (api *LolesportsAPI) GetStandings(hl string, tournamentID string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getStandings", map[string]string{"hl": hl, "tournamentId": tournamentID})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetSchedule retrieves schedule
func (api *LolesportsAPI) GetSchedule(hl string, leagueID string, pageToken string) (map[string]interface{}, error) {
	params := map[string]string{"hl": hl}
	if leagueID != "" {
		params["leagueId"] = leagueID
	}
	if pageToken != "" {
		params["pageToken"] = pageToken
	}

	response, err := api.getRequest(apiURL+"/getSchedule", params)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetLive retrieves live data
func (api *LolesportsAPI) GetLive(hl string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getLive", map[string]string{"hl": hl})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetCompletedEvents retrieves completed events
func (api *LolesportsAPI) GetCompletedEvents(hl string, tournamentID string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getCompletedEvents", map[string]string{"hl": hl, "tournamentId": tournamentID})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetEventDetails retrieves event details
func (api *LolesportsAPI) GetEventDetails(matchID string, hl string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getEventDetails", map[string]string{"hl": hl, "id": matchID})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetGames retrieves games data
func (api *LolesportsAPI) GetGames(hl string, matchID string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getGames", map[string]string{"hl": hl, "id": matchID})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetTeams retrieves teams data
func (api *LolesportsAPI) GetTeams(hl string, teamSlug string) (map[string]interface{}, error) {
	response, err := api.getRequest(apiURL+"/getTeams", map[string]string{"hl": hl, "id": teamSlug})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetWindow retrieves window data
func (api *LolesportsAPI) GetWindow(gameID string, startingTime string) (map[string]interface{}, error) {
	response, err := api.getRequest(liveStatsAPI+"/window/"+gameID, map[string]string{"startingTime": startingTime})
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetDetails retrieves details data
func (api *LolesportsAPI) GetDetails(gameID string, startingTime string, participantIDs []string) (map[string]interface{}, error) {
	params := map[string]string{"startingTime": startingTime}
	if len(participantIDs) > 0 {
		params["participantIds"] = fmt.Sprintf("%v", participantIDs)
	}

	response, err := api.getRequest(liveStatsAPI+"/details/"+gameID, params)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(response.Data, &data); err != nil {
		return nil, err
	}
	return data, nil
}
