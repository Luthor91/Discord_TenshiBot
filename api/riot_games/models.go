package riot_games

// ChampionMastery représente les données de maîtrise d'un champion
type ChampionMastery struct {
	ChampionID     int `json:"championId"`
	ChampionLevel  int `json:"championLevel"`
	ChampionPoints int `json:"championPoints"`
}

// ChampionData représente les données d'un champion
type ChampionData struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// SummonerProfile représente la structure des données d'un profil d'invocateur
type SummonerProfile struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountId"`
	PUUID       string `json:"puuid"`
	Name        string `json:"name"`
	TagLine     string `json:"tagLine"`
	ProfileIcon int    `json:"profileIconId"`
}

// EsportMatch définit la structure pour les informations d'un match eSports
type EsportMatch struct {
	MatchID   string
	StartTime string
	Team1     string
	Team2     string
	League    string
	// Ajoutez d'autres champs nécessaires ici
}

type LobbyEvent struct {
	ID        string `json:"id"`
	StartTime string `json:"startTime"`
	Team1     string `json:"team1"`
	Team2     string `json:"team2"`
	League    string `json:"league"`
}
