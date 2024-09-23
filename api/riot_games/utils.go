package riot_games

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/Luthor91/Tenshi/config"
)

// getStatic fetches static data from the specified URL and decodes it into the provided target structure.
func getStatic(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data from URL %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	err = json.Unmarshal(body, &target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return nil
}

// GetChampionsNameByIds convertit une liste d'IDs de champions en leurs noms
func GetChampionsNameByIds(ids []float64) ([]string, error) {
	// URL pour les données des champions
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/fr_FR/champion.json", config.AppConfig.LoLPatchVersion)

	// Récupération des données des champions
	var result struct {
		Data map[string]ChampionData `json:"data"`
	}
	err := getStatic(url, &result)
	if err != nil {
		return nil, err
	}

	// Création de la map pour les noms des champions
	championNames := make(map[int]string)
	for _, champ := range result.Data {
		championID, _ := strconv.Atoi(champ.Key)
		championNames[championID] = champ.Name
	}

	// Création d'une liste pour les noms des champions demandés
	var names []string
	for _, id := range ids {
		idInt := int(id)
		if name, exists := championNames[idInt]; exists {
			names = append(names, name)
		} else {
			names = append(names, "Unknown")
		}
	}

	return names, nil
}

// GetChampionNameById convertit un ID de champion en son nom
func GetChampionNameById(id int) (string, error) {
	// URL pour les données des champions
	url := fmt.Sprintf("http://ddragon.leagueoflegends.com/cdn/%s/data/fr_FR/champion.json", config.AppConfig.LoLPatchVersion)

	// Récupération des données des champions
	var result struct {
		Data map[string]ChampionData `json:"data"`
	}
	err := getStatic(url, &result)
	if err != nil {
		return "", err
	}

	// Recherche du nom du champion correspondant à l'ID
	championID := strconv.Itoa(id) // Conversion de l'ID en string pour correspondre à la clé des champions
	for _, champ := range result.Data {
		if champ.Key == championID {
			return champ.Name, nil
		}
	}

	// Si aucun champion n'a été trouvé pour l'ID donné
	return "Unknown", nil
}

// GetSummonerProfile récupère les informations du profil d'un invocateur
func GetSummonerProfile(name, tag string) (interface{}, error) {
	// Créer l'URL de la requête pour obtenir le profil d'invocateur
	url := fmt.Sprintf("%s/riot/account/v1/accounts/by-riot-id/%s/%s", config.AppConfig.RiotBaseURL, name, tag)

	// Préparer la requête HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Ajouter la clé API dans les en-têtes
	req.Header.Add("X-Riot-Token", config.AppConfig.RiotAPIKey)

	// Envoyer la requête
	client := config.AppConfig.Client
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Vérifier le code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get summoner profile, status code: %d", resp.StatusCode)
	}

	// Décoder la réponse JSON
	var profile Summoner
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// CalculateWinrate calcule le taux de victoire en pourcentage
func CalculateWinrate(wins int, totalGames int) float64 {
	if totalGames == 0 {
		return 0 // Eviter la division par zéro, retour d'un taux de victoire de 0%
	}
	winrate := (float64(wins) / float64(totalGames)) * 100
	return winrate
}

// FormatWinrate formatte le taux de victoire en pourcentage avec deux décimales
func FormatWinrate(winrate float64) string {
	return fmt.Sprintf("%.2f%%", winrate)
}
