package riot_games

// ChampionDataResponse représente la réponse globale de l'API pour les champions
type ChampionDataResponse struct {
	Data map[string]ChampionDataExtended `json:"data"`
}

// ChampionMastery représente les données de maîtrise d'un champion
type ChampionMastery struct {
	ChampionID     int `json:"championId"`
	ChampionLevel  int `json:"championLevel"`
	ChampionPoints int `json:"championPoints"`
}

// ChampionData représente les données détaillées d'un champion
type ChampionData struct {
	ID    string `json:"id"`    // Identifiant unique du champion
	Key   string `json:"key"`   // Identifiant numérique du champion sous forme de chaîne
	Name  string `json:"name"`  // Nom du champion
	Title string `json:"title"` // Titre du champion
	Blurb string `json:"blurb"` // Description courte du champion
}

// ChampionDataExtended représente les données étendues d'un champion
type ChampionDataExtended struct {
	ID    string `json:"id"`    // Identifiant unique du champion
	Key   string `json:"key"`   // Clé unique numérique en tant que string
	Name  string `json:"name"`  // Nom du champion
	Title string `json:"title"` // Titre du champion
	Blurb string `json:"blurb"` // Description ou résumé du champion
	// Vous pouvez ajouter d'autres champs si nécessaire, comme:
	Tags    []string `json:"tags"`    // Catégories ou rôles du champion (ex: Fighter, Mage)
	Partype string   `json:"partype"` // Type de ressource utilisé (ex: Mana, Energy)
	Info    Info     `json:"info"`    // Détails supplémentaires sur les stats du champion (difficulté, attaque, etc.)
	Image   Image    `json:"image"`   // Détails de l'image associée au champion
	Stats   Stats    `json:"stats"`   // Statistiques du champion (HP, AD, Armor, etc.)
}

// Info fournit des informations sur les attributs du champion
type Info struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

// Image contient les détails de l'image associée au champion
type Image struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

// Stats représente les statistiques du champion
type Stats struct {
	HP                   float64 `json:"hp"`
	HPPerLevel           float64 `json:"hpperlevel"`
	MP                   float64 `json:"mp"`
	MPPerLevel           float64 `json:"mpperlevel"`
	MoveSpeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	ArmorPerLevel        float64 `json:"armorperlevel"`
	SpellBlock           float64 `json:"spellblock"`
	SpellBlockPerLevel   float64 `json:"spellblockperlevel"`
	AttackRange          float64 `json:"attackrange"`
	HPRegen              float64 `json:"hpregen"`
	HPRegenPerLevel      float64 `json:"hpregenperlevel"`
	MPRegen              float64 `json:"mpregen"`
	MPRegenPerLevel      float64 `json:"mpregenperlevel"`
	Crit                 float64 `json:"crit"`
	CritPerLevel         float64 `json:"critperlevel"`
	AttackDamage         float64 `json:"attackdamage"`
	AttackDamagePerLevel float64 `json:"attackdamageperlevel"`
	AttackSpeedPerLevel  float64 `json:"attackspeedperlevel"`
	AttackSpeed          float64 `json:"attackspeed"`
}

// ChampionBasicInfo représente les données de base d'un champion
type ChampionBasicInfo struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Summoner struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
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
